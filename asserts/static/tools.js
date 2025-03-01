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