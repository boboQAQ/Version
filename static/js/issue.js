$(function() {

    var Data
    var num = 0;
    //创建websocket链接
    var socket = new WebSocket("ws://127.0.0.1:8080/WS3");

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

    //给发布列表动态展示函数
    function showTable(data) {
         //向表格动态添加数据
      var versiontab = $('#versiontable');
      list = data.servicelist;
      for(var i = 0; i < list.length; i++) {
        //if(list[i].servicenumber == "&&&")continue;
        //先添加大版本的创建时间和发布时间备注等，之后再更改为list的创建时间的message等
        versiontab.append('<tr class="success"> ' +
        '<td>' + i + '</td>' + 
        '<td>' + data.versionnumber + '</td>' +
        '<td>' +  list[i].servicename +  '</td>' +
        '<td>' +  list[i].servicenumber +  '</td>' +
        '<td>' + format(data.issuetime) + '</td>' +
        '<td>' + format(data.creattime) + '</td>' +
        '<td>' + 
        '<button id="button1" type="button" class="but" value="0">合并</button>' +
        '<button id="button2" type="button" class="but" value="1">发布</button>' + 
        '</td>' +
        '</tr>')
      }
    }
    $(document).on('click','#button1',function(){
       
        console.log("点击合并");
        var send = $(this).parents("tr").find('td').eq(3).text() + document.getElementById("button1").value;
        console.log(send);
        socket.send(send);

    })
    $(document).on('click','#button2',function(){
       
        console.log("点击发布");
        console.log($("#slpk1").val());
        console.log( $(this).parents("tr").find('td').eq(2).text());
        var send = $(this).parents("tr").find('td').eq(3).text() + document.getElementById("button2").value;
        send = send + " " + $("#slpk1").val();
        socket.send(send);
        refresh();

    })

 
    socket.onmessage = function(event) {
        num++;
        if(num>1)
        {
            var data = JSON.parse(event.data);
            window.alert(data);
            return ;
        }
        //解析json，之后初始化加载的更新页面
        var data = JSON.parse(event.data);
            console.log(data)
            var data1 = data.services;
            var data2 = data.versions;
            console.log(data1)
            var select2 = $("#slpk2");       //下面给服务列表动态添加服务 
            for(var i = 0; i < data1.length; i++) {
                var str1 =  data1[i].serviceid;
                var str2 = data1[i].servicename ;
                select2.append("<option value='"+str1+"'>"+str2+"</option>"); 
            }
            select2.selectpicker('refresh');

            Data = data.versions;              //将版本信息赋值给全局变量
            console.log("revice:", data2);     //输出解析之后的后台文件
            var select = $("#slpk1");        //给下拉框定义别名
            var list = data2[0].servicelist //默认选择了第一个，所以这是它的服务列表
            showTable(data2[0]);
            var i = 1;
            for(var j = 0; j < data2.length; j++){   //使用jQuery动态给下拉框添加option
                if(data2[j].versionnumber == "")data2[j].versionnumber = "无效版本"+i++;
                select.append("<option value='"+data2[j].id+"'>"+data2[j].versionnumber+"</option>"); 
            }
            select.selectpicker('refresh');   //刷新下拉框
        
            for(var j = 0; j < list.length; j++) {   //根据后台数据，在服务列表的下拉框中选出该大版本下挂钩的服务
                var str
                str =  list[j].servicenumber;
                $("#slpk2 option[value='"+str+"']").prop("selected","selected");  
            }
            select2.selectpicker('refresh');
            document.getElementById("comment").value = data2[0].comment; //将选中的版本号的备注显示在页面
        
    }


    $('#slpk1').change( function() {
        var data = Data;
        var id = $("#slpk1").val();
        for(var j = 0; j < data.length; j++) {
            if(data[j].id == id) {
                //清空表格
                $("#versiontable tbody").html("");
                showTable(data[j]);
                var list = data[j].servicelist;
                $("#slpk2").find("option:selected").attr("selected", false);
                for(var k = 0; k < list.length; k++) {
                    var str
                    str = list[k].servicenumber;
                    $("#slpk2 option[value='"+str+"']").prop("selected","selected");       
                }
                $("#slpk2").selectpicker('refresh');
                document.getElementById("comment").value = data[j].comment;
                break;
            }
        }
    });

    $(window).keydown(function (event) {
        // 13是回车的键位
        if (event.which === 13) {
            sendMessage();
            typing = false;
        }
    });

    // 发送按钮点击事件
    $("#send1").click(function () {
        sendMessage();
    });
    function sendMessage() {
        console.log("点击发布")
        var t = $('#form1').serializeArray();
        // for(var j = 0; j < Data.length; j++){
        //     if(Data[j].id == t[0].value){
        //         t[0].value = Data[j].versionnumber;
        //         break;
        //     }
        // }
        console.log(JSON.stringify(t));
        var json_str = JSON.stringify(t);
        
        socket.send(json_str);
        

    }

});