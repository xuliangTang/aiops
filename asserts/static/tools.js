//判断 是否是table格式
function is_table(obj){
    let s=0;
    for(let key in obj){
        if(key==="Cols"){
            s++
        }
        if(key==="Rows"){

            s++
        }
    }

    return s===2;
}

function render_table(obj){
    const table = document.createElement("table");
    let tr = table.insertRow(-1);                   // table row.

    for (let i = 0; i < obj.Cols.length; i++) {
        let th = document.createElement("th");      // table header.
        th.innerHTML = obj.Cols[i];
        tr.appendChild(th);
    }

    for (let i = 0; i < obj.Rows.length; i++) {
        tr = table.insertRow(-1);
        for (let j = 0; j < obj.Cols.length; j++) {
            let tabCell = tr.insertCell(-1);
            tabCell.innerHTML = obj.Rows[i][obj.Cols[j]];
        }
    }

    table.className="gridtable" //设置演示
    const divResult = document.getElementById('result');
    divResult.innerHTML = "";

    divResult.appendChild(table);

}

// 判断 是否是shell 模式
function if_shell_return_url(obj){
    if (obj.startsWith("shell:")) {
        return  obj.substring("shell:".length)
    }
    return '';
}

function str2utf8(str) {
    let encoder = new TextEncoder('utf8');
    return encoder.encode(str);
}
function ab2str(buf) {
    let encodedString = String.fromCodePoint.apply(null, new Uint8Array(buf));
    encodedString=escape(encodedString)
    return    decodeURIComponent(encodedString);//没有这一步中文会乱码
}

function createTerminalDiv(url){
    const divResult = document.getElementById('result');
    let divTerm=document.createElement("div");
    divTerm.id=md5(url)
    divResult.appendChild(divTerm);
    createTerminal(divTerm,url)
}

function createTerminal(container,url){
    console.log("url是"+url)
    var websocket = new WebSocket("ws://" + window.location.hostname + ":" + window.location.port + url);
    websocket.binaryType = "arraybuffer";
    websocket.onopen = function(evt) {
        const term =  new Terminal({
            rendererType: 'canvas',
            screenKeys: true,
            useStyle: true,
            cursorBlink: true,
            wraparound:false,
            rows: 14, //行数
            cols: 80, // 不指定行数，自动回车后光标从下一行开始
            fontSize:16,
            allowMouse: true,
            windowsMode: true,
            theme: {
                foreground: "#7e9192", //字体
                background: "#000", //背景色
                cursor: "help", //设置光标
                lineHeight: 22,
                paddingTop:"20px",
                marginTop:20
            }
        });


        term.onData((data)=> {
            websocket.send(str2utf8(data));
        });
        term.prompt = () => {
            term.writeln("\n 这是容器终端");
        };
        term.prompt();
        term.open(container);
        websocket.onmessage=function(e){
            term.write(ab2str(e.data));
        };
        websocket.onclose = function(e){
            console.log("close");
        }
        websocket.onerror = function(e){
            console.log(e);
        }
    }

    let divTerm=document.createElement("div");
    document.createElement("div").appendChild(divTerm);
}