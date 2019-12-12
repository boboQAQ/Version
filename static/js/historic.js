$(function() {

    var Data
    //创建websocket链接
    var socket = new WebSocket("ws://127.0.0.1:8080/WS4");

    socket.onopen = function() {
        console.log("websocket open");
        connected = true;
    };

    socket.onclose = function() {
        console.log("websocket close");
        connected = false;
    };
    //将不足10的补齐0
    function getzf(num) {  
        if(parseInt(num) < 10){  
            num = '0'+num;  
        }  
        return num;  
    }
    //将此类2019-12-06T09:40:11Z 格式的时间转换为正常时间
    function format(date1) {
        var date = new Date(date1 );
       date = date -  8 * 3600 * 1000;
       var d = new Date(date);
        var times= getzf(d.getFullYear()) + '-' + getzf((d.getMonth() + 1)) + '-'
        + getzf(d.getDate()) + ' ' + getzf(d.getHours()) + ':' + getzf(d.getMinutes())
        + ':' + getzf(d.getSeconds());

         return times;
    }
   
    socket.onmessage = function(event) {
       //解析json，之后初始化加载的更新页面
       var data = JSON.parse(event.data);
       var data1 = data.services;
       var data2 = data.versions;

       var select2 = $("#slpk2");       //下面给服务列表动态添加服务 
       for(var i = 0; i < data1.length; i++) {
           var str1 = data1[i].servicename + "&" + data1[i].servicenumber;
           var str2 = data1[i].servicename + "服务版本号" + data1[i].servicenumber;
           select2.append("<option value='"+str1+"'>"+str2+"</option>"); 
       }
      select2.selectpicker('refresh');
      //向表格动态添加数据
      console.log(data2);
      var versiontab = $('#versiontable');
      for(var i = 0; i < data2.length; i++) {

        var list = data2[i].servicelist;
        var str1 = "";
        for(var j = 0; j < list.length; j++) {   //根据后台数据，在服务列表的下拉框中选出该大版本下挂钩的服务
            str1 += "&lt" + list[j].servicename + " " + list[j].servicenumber + "&gt";
       }
       console.log(str1);

        versiontab.append('<tr class="success"> ' +
		'<td>' + i + '</td>' + 
		'<td>' + data2[i].versionnumber + '</td>' +
        '<td>' +  str1+  '</td>' +
        '<td>' + format(data2[i].issuetime) + '</td>' +
        '<td>' + format(data2[i].creattime) + '</td>' +
        '<td>' + data2[i].comment + '</td>' +
        '</tr>')
      }


       Data = data.versions;              //将版本信息赋值给全局变量
       console.log("revice:", data2);     //输出解析之后的后台文件
       var select = $("#slpk1");        //给下拉框定义别名
       var list = data2[0].servicelist //默认选择了第一个，所以这是它的服务列表
       var i = 1;
       for(var j = 0; j < data2.length; j++){   //使用jQuery动态给下拉框添加option
           if(data2[j].versionnumber == "")data2[j].versionnumber = "无效版本"+i++;
           select.append("<option value='"+data2[j].id+"'>"+data2[j].versionnumber+"</option>"); 
       }
      select.selectpicker('refresh');   //刷新下拉框
      
      for(var j = 0; j < list.length; j++) {   //根据后台数据，在服务列表的下拉框中选出该大版本下挂钩的服务
           var str
           str = list[j].servicename + "&" + list[j].servicenumber;
           $("#slpk2 option[value='"+str+"']").prop("selected","selected");  
      }
      select2.selectpicker('refresh');
      document.getElementById("comment").value = data2[0].comment; //将选中的版本号的备注显示在页面
    }

    //实现点击表头表格自动排序（含数字、字符串、日期）
    var tbody = document.querySelector('#versiontable').tBodies[0];
    var th = document.querySelector('#versiontable').tHead.rows[0].cells;
    var td = tbody.rows;
    for (var i = 0;i < th.length;i++){
        th[i].flag = 1;
        th[i].onclick = function(){
            sort(this.getAttribute('data-type'),this.flag,this.cellIndex);
            this.flag = -this.flag;
        };
    };
    function sort(str,flag,n){
        var arr = []; //存放DOM
        for (var i = 0;i < td.length;i++){
            arr.push(td[i]);
        };
        //排序
        arr.sort(function(a,b){
            return method(str,a.cells[n].innerHTML,b.cells[n].innerHTML) * flag;
        });
        //添加
        for (var i = 0;i < arr.length;i++){
            tbody.appendChild(arr[i]);
        };
    };
    //排序方法
    function method(str,a,b){
        switch (str){
        case 'num': //数字排序
            return a-b;
            break;
        case 'string': //字符串排序
            return a.localeCompare(b);
            break;
        default:  //日期排序，IE8下'2012-12-12'这种格式无法设置时间，替换成'/'
            return new Date(a.split('-').join('/')).getTime()-new Date(b.split('-').join('/')).getTime();
        };
    };


    $('#slpk1').change( function() {
        var data = Data;
        var id = $("#slpk1").val();
        for(var j = 0; j < data.length; j++) {
            if(data[j].id == id) {
                var list = data[j].servicelist;
                $("#slpk2").find("option:selected").attr("selected", false);
                for(var k = 0; k < list.length; k++) {
                    var str
                    str = list[k].servicename + "&" + list[k].servicenumber;
                    $("#slpk2 option[value='"+str+"']").prop("selected","selected");       
                }
                $("#slpk2").selectpicker('refresh');
                document.getElementById("comment").value = data[j].comment;
                break;
            }
        }
    });

   

});