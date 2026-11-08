package stringarray

import (
	"encoding/base64"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/tdewolff/parse/v2/js"

	"veiljs/internal/jsstr"
	"veiljs/internal/names"
	"veiljs/internal/rng"
	"veiljs/parser"
)

type Options struct {
	Encoding  string
	Threshold float64
}

func Apply(ast *js.AST, r *rng.RNG, g *names.Gen, o Options) {
	skip := collectSkip(ast)
	arr := g.Next()
	dec := g.Next()

	key := ""
	if o.Encoding == "rc4" {
		key = randomKey(r)
	}

	assign := map[*js.LiteralExpr]int{}
	seen := map[string]int{}
	var entries []string

	parser.Rewrite(ast, func(n js.INode) js.INode {
		lit, ok := n.(*js.LiteralExpr)
		if !ok || lit.TokenType != js.StringToken || skip[lit] {
			return n
		}
		s := jsstr.Unquote(lit.Data)
		if i, ok := seen[s]; ok {
			if 0 <= i {
				assign[lit] = i
			}
			return n
		}
		if r.Float64() >= o.Threshold {
			seen[s] = -1
			return n
		}
		i := len(entries)
		seen[s] = i
		assign[lit] = i
		entries = append(entries, string(jsstr.Quote(encode(o.Encoding, s, key))))
		return n
	})

	if len(entries) == 0 {
		return
	}

	src := buildSource(arr, dec, o.Encoding, key, entries, g)
	mod, err := parser.Parse(src)
	if err != nil {
		return
	}

	decVar := &js.Var{Data: []byte(dec)}
	parser.Rewrite(ast, func(n js.INode) js.INode {
		if lit, ok := n.(*js.LiteralExpr); ok {
			if i, ok := assign[lit]; ok {
				return lookup(decVar, i)
			}
		}
		return n
	})

	ast.List = splice(ast.List, mod.List)
}

func lookup(dec *js.Var, i int) js.IExpr {
	return &js.CallExpr{
		X:    dec,
		Args: js.Args{List: []js.Arg{{Value: &js.LiteralExpr{TokenType: js.IntegerToken, Data: []byte(strconv.Itoa(i))}}}},
	}
}

func splice(list, inject []js.IStmt) []js.IStmt {
	i := 0
	for i < len(list) {
		if _, ok := list[i].(*js.Comment); ok {
			i++
			continue
		}
		if _, ok := list[i].(*js.DirectivePrologueStmt); ok {
			i++
			continue
		}
		break
	}
	out := make([]js.IStmt, 0, len(list)+len(inject))
	out = append(out, list[:i]...)
	out = append(out, inject...)
	out = append(out, list[i:]...)
	return out
}

func encode(enc, s, key string) string {
	b := []byte(s)
	switch enc {
	case "base64":
		return base64.StdEncoding.EncodeToString(b)
	case "hex":
		return hex.EncodeToString(b)
	case "rc4":
		return base64.StdEncoding.EncodeToString(rc4(b, []byte(key)))
	default:
		return s
	}
}

func rc4(data, key []byte) []byte {
	var s [256]byte
	for i := range s {
		s[i] = byte(i)
	}
	j := 0
	for i := 0; i < 256; i++ {
		j = (j + int(s[i]) + int(key[i%len(key)])) & 255
		s[i], s[j] = s[j], s[i]
	}
	out := make([]byte, len(data))
	i := 0
	j = 0
	for k, b := range data {
		i = (i + 1) & 255
		j = (j + int(s[i])) & 255
		s[i], s[j] = s[j], s[i]
		out[k] = b ^ s[(int(s[i])+int(s[j]))&255]
	}
	return out
}

func randomKey(r *rng.RNG) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = alphabet[r.Intn(len(alphabet))]
	}
	return string(b)
}

func buildSource(arr, dec, enc, key string, entries []string, g *names.Gen) string {
	var sb strings.Builder
	sb.WriteString("var ")
	sb.WriteString(arr)
	sb.WriteString("=[")
	for i, e := range entries {
		if i != 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(e)
	}
	sb.WriteString("];")
	switch enc {
	case "hex":
		h := g.Next()
		u := g.Next()
		sb.WriteString(hexSrc(h))
		sb.WriteString(utf8Src(u))
		sb.WriteString("function " + dec + "(i){return " + u + "(" + h + "(" + arr + "[i]));}")
	case "base64":
		b := g.Next()
		u := g.Next()
		sb.WriteString(b64Src(b))
		sb.WriteString(utf8Src(u))
		sb.WriteString("function " + dec + "(i){return " + u + "(" + b + "(" + arr + "[i]));}")
	case "rc4":
		b := g.Next()
		c := g.Next()
		u := g.Next()
		sb.WriteString(b64Src(b))
		sb.WriteString(rc4Src(c))
		sb.WriteString(utf8Src(u))
		sb.WriteString("function " + dec + "(i){return " + u + "(" + c + "(" + b + "(" + arr + "[i])," + string(jsstr.Quote(key)) + "));}")
	default:
		sb.WriteString("function " + dec + "(i){return " + arr + "[i];}")
	}
	return sb.String()
}

func hexSrc(name string) string {
	return "function " + name + "(s){var o='',i;for(i=0;i<s.length;i+=2){o+=String.fromCharCode(parseInt(s.substring(i,i+2),16))}return o}"
}

func b64Src(name string) string {
	return "function " + name + "(s){var t='ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/',o='',n=0,b=0,i,c;for(i=0;i<s.length;i++){c=t.indexOf(s.charAt(i));if(c<0){continue}b=(b<<6)|c;n+=6;if(n>=8){n-=8;o+=String.fromCharCode((b>>n)&255)}}return o}"
}

func rc4Src(name string) string {
	return "function " + name + "(s,k){var a=[],i,j=0,x,o='';for(i=0;i<256;i++){a[i]=i}for(i=0;i<256;i++){j=(j+a[i]+k.charCodeAt(i%k.length))&255;x=a[i];a[i]=a[j];a[j]=x}i=0;j=0;for(var y=0;y<s.length;y++){i=(i+1)&255;j=(j+a[i])&255;x=a[i];a[i]=a[j];a[j]=x;o+=String.fromCharCode(s.charCodeAt(y)^a[(a[i]+a[j])&255])}return o}"
}

func utf8Src(name string) string {
	return "function " + name + "(s){var o='',i=0,c,u;while(i<s.length){c=s.charCodeAt(i++);if(c<128){o+=String.fromCharCode(c)}else if(c<224){o+=String.fromCharCode(((c&31)<<6)|(s.charCodeAt(i++)&63))}else if(c<240){o+=String.fromCharCode(((c&15)<<12)|((s.charCodeAt(i++)&63)<<6)|(s.charCodeAt(i++)&63))}else{u=((c&7)<<18)|((s.charCodeAt(i++)&63)<<12)|((s.charCodeAt(i++)&63)<<6)|(s.charCodeAt(i++)&63);u-=65536;o+=String.fromCharCode(55296+(u>>10),56320+(u&1023))}}return o}"
}

func collectSkip(root js.INode) map[*js.LiteralExpr]bool {
	s := &skipper{skip: map[*js.LiteralExpr]bool{}}
	js.Walk(s, root)
	return s.skip
}

type skipper struct {
	skip  map[*js.LiteralExpr]bool
	depth int
}

func (s *skipper) Enter(n js.INode) js.IVisitor {
	switch n := n.(type) {
	case *js.CallExpr:
		if isModuleCall(n) {
			s.depth++
		}
	case *js.LiteralExpr:
		if 0 < s.depth && n.TokenType == js.StringToken {
			s.skip[n] = true
		}
	}
	return s
}

func (s *skipper) Exit(n js.INode) {
	if c, ok := n.(*js.CallExpr); ok && isModuleCall(c) {
		s.depth--
	}
}

func isModuleCall(c *js.CallExpr) bool {
	switch x := c.X.(type) {
	case *js.Var:
		return string(x.Name()) == "require"
	case *js.LiteralExpr:
		return x.TokenType == js.ImportToken
	}
	return false
}
