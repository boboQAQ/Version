$(function() {

    //作为全局变量记录从后台得到的数据
    var Data 

    //====================webSocket连接======================
    // 创建一个webSocket连接
    var socket = new WebSocket("ws://127.0.0.1:8080/WS2");

    // 当webSocket连接成功的回调函数
    socket.onopen = function () {
        console.log("webSocket open");
        connected = true;
    };

    // 断开webSocket连接的回调函数
    socket.onclose = function () {
        console.log("webSocket close");
        connected = false;
    };


    //接收到后台查询到的信息
    socket.onmessage = function(event) {
        //解析json，之后初始化加载的更新页面
        var data = JSON.parse(event.data);
        var data1 = data.services;
        var data2 = data.versions;

        var select2 = $("#slpk2");       //下面给服务列表动态添加服务 
        for(var i = 0; i < data1.length; i++) {
            var str1 = data1[i].servicename + "&" + data1[i].servicenumber;
            var str2 = data1[i].servicename ;
            select2.append("<option value='"+str1+"'>"+str2+"</option>"); 
        }
       select2.selectpicker('refresh');


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
    //使用change监听版本号的下拉框option的改变，去更新服务列表
    $('#slpk1').change( function()  {
        // console.log('change');
        // console.log($('#cont select option:selected').text());
        var data = Data;
        console.log(data);
        var jobtitle = $("#slpk1").val();  //取出选中的option的value，也就是它们在数据库中的ID，之后根据ID去选着服务更新服务列表
        console.log(jobtitle);
        for(var j = 0; j < data.length; j++){
            if(data[j].id == jobtitle) {
                var list = data[j].servicelist;
                $('#slpk2').find("option:selected").attr("selected", false);
                for(var k = 0; k < list.length; k++) {
                    var str
                    str = list[k].servicename + "&" + list[k].servicenumber;
                    $("#slpk2 option[value='"+str+"']").prop("selected","selected");    //       
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
        console.log("进入修改")
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
        $('#slpk2').find("option:selected").attr("selected", false);
        $('#slpk2').selectpicker('refresh'); 

    }

});