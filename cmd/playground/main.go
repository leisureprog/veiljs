package main

import (
	"encoding/json"
	"net/http"
	"os"

	"veiljs"
)

const page = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>veiljs playground</title>
<style>
:root{color-scheme:dark}
*{margin:0;padding:0;box-sizing:border-box}
body{background:#0b1220;color:#cbd5e1;font:15px/1.5 -apple-system,"Segoe UI",sans-serif;height:100vh;display:flex;flex-direction:column}
header{display:flex;align-items:center;gap:10px;padding:14px 24px;border-bottom:1px solid #1e293b}
h1{font-size:18px;font-weight:600;color:#f1f5f9}
h1 span{color:#19d3c5}
.controls{display:flex;flex-wrap:wrap;gap:14px;align-items:center;padding:10px 24px;border-bottom:1px solid #1e293b;font-size:13px;color:#94a3b8}
label{display:flex;align-items:center;gap:6px;cursor:pointer;user-select:none}
input[type=checkbox]{accent-color:#19d3c5}
select{background:#0f172a;color:#cbd5e1;border:1px solid #1e293b;border-radius:6px;padding:4px 8px;outline:0}
button{margin-left:auto;background:#19d3c5;color:#0b1220;border:0;border-radius:8px;padding:8px 22px;font-weight:600;font-size:14px;cursor:pointer}
button:active{transform:translateY(1px)}
main{flex:1;display:grid;grid-template-columns:1fr 1fr;gap:1px;background:#1e293b;min-height:0}
section{display:flex;flex-direction:column;background:#0f172a;min-height:0}
.bar{padding:8px 14px;font:12px Menlo,monospace;color:#64748b;border-bottom:1px solid #1e293b;display:flex;justify-content:space-between}
textarea{flex:1;background:transparent;border:0;outline:0;resize:none;padding:14px;font:13px/1.6 Menlo,monospace;color:#cbd5e1;white-space:pre}
#out{color:#94a3b8}
</style>
</head>
<body>
<header><h1>veil<span>js</span> playground</h1></header>
<div class="controls">
<label><input type="checkbox" id="sa" checked>string array</label>
<select id="enc"><option>base64</option><option>none</option><option>hex</option><option>rc4</option></select>
<label><input type="checkbox" id="cff" checked>control flow</label>
<label><input type="checkbox" id="num">numbers</label>
<label><input type="checkbox" id="sim" checked>simplify</label>
<label><input type="checkbox" id="uni">unicode</label>
<label><input type="checkbox" id="ren">rename</label>
<button id="run">Obfuscate &#8594;</button>
</div>
<main>
<section>
<div class="bar"><span>input.js</span><span id="inSize"></span></div>
<textarea id="in" spellcheck="false">function hello(name) {
  var msg = "Hello, " + name + "!";
  var key = "s3cr3t-key";
  for (var i = 0; i < 3; i++) {
    console.log(msg, key, i);
  }
  return 42;
}
hello("world");</textarea>
</section>
<section>
<div class="bar"><span>obfuscated.js</span><span id="outSize"></span></div>
<textarea id="out" spellcheck="false" readonly></textarea>
</section>
</main>
<script>
var $ = function(id) { return document.getElementById(id); };
function run() {
	var body = {
		Code: $("in").value,
		Compact: true,
		StringArray: $("sa").checked,
		StringArrayEncoding: [$("enc").value],
		StringArrayThreshold: 1,
		ControlFlowFlattening: $("cff").checked,
		ControlFlowFlatteningThreshold: 1,
		NumbersToExpressions: $("num").checked,
		Simplify: $("sim").checked,
		UnicodeEscapeSequence: $("uni").checked,
		RenameGlobals: $("ren").checked
	};
	$("inSize").textContent = $("in").value.length + " B";
	fetch("/api/obfuscate", { method: "POST", body: JSON.stringify(body) })
		.then(function(r) { return r.json(); })
		.then(function(d) {
			$("out").value = d.Output;
			$("outSize").textContent = d.Output.length + " B";
		});
}
$("run").onclick = run;
document.addEventListener("keydown", function(e) {
	if ((e.metaKey || e.ctrlKey) && e.key === "Enter") run();
});
run();
</script>
</body>
</html>`

type request struct {
	Code string
	veiljs.Config
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(page))
	})
	http.HandleFunc("/api/obfuscate", obfuscate)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}

func obfuscate(w http.ResponseWriter, r *http.Request) {
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	out, err := veiljs.ObfuscateWithConfig([]byte(req.Code), req.Config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct{ Output string }{out})
}
