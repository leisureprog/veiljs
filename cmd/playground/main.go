package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/leisureprog/veiljs"
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
.gh{margin-left:auto;color:#64748b;display:flex}
.gh:hover{color:#19d3c5}
main{flex:1;display:grid;grid-template-columns:1fr 1fr;gap:1px;background:#1e293b;min-height:0}
section{display:flex;flex-direction:column;background:#0f172a;min-height:0}
.bar{padding:8px 14px;font:12px Menlo,monospace;color:#64748b;border-bottom:1px solid #1e293b;display:flex;justify-content:space-between}
textarea{flex:1;background:transparent;border:0;outline:0;resize:none;padding:14px;font:13px/1.6 Menlo,monospace;color:#cbd5e1;white-space:pre}
#out{color:#94a3b8}
</style>
</head>
<body>
<header><h1>veil<span>js</span> playground</h1><a class="gh" href="https://github.com/leisureprog/veiljs" target="_blank" rel="noopener"><svg width="22" height="22" viewBox="0 0 16 16" fill="currentColor"><path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/></svg></a></header>
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
