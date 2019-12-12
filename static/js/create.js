
 $(function() {
    $('.selectpicker').selectpicker({
    'selectedText':'cat',
     'noneSelectedText':'请选择',
     'deselectAllText':'全不选',
     'selectAllText': '全选',
 })
$(window).on('load', function() { 

    $('.selectpicker').selectpicker('refresh'); 
  
}); 
//实时监听版本号的input框
 $("#inputMessage1").bind("input propertychange",function(event){
    console.log($("#inputMessage1").val())
});
//实时监听备注的input框
$("#inputMessage2").bind("input propertychange",function(event){
    console.log($("#inputMessage2").val())
});
//====================webSocket连接======================
    // 创建一个webSocket连接
    var socket = new WebSocket("ws://127.0.0.1:8080/WS1");

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
    socket.onmessage = function (event) {
        var select = $("#slpk"); 
        var data = JSON.parse(event.data);
        for(var i = 0; i < data.length; i++) {
            var str1 = data[i].servicename + "&" + data[i].servicenumber;
            var str2 = data[i].servicename + "服务版本号" + data[i].servicenumber;
            select.append("<option value='"+str1+"'>"+str2+"</option>"); 
        }
        select.selectpicker('refresh');

    }



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
        var t = $('#form1').serializeArray();
        //var t = $('#form1').serializeArray(); 
        // 默认是json 格式，将表单序列化
        //传统的for循环,将值打印出来                        
        // console.log(t.length);
        // for (var i = 0;i<t.length;i++) {
            
        //     console.info(t[i]);
        
        //     d[t[i].name] =  t[i].value;
        // }
    
        //jQuery的循环
        /*
        $.each(t, function() {
        
        //console.info(t)
        
        d[this.name] = this.value;
        
        console.info(d)
        
        }); 
        */
        console.info(JSON.stringify(t));   //从json对象解析出json字符串
        
        var json_str = JSON.stringify(t);
        
        socket.send(json_str);
        $('#slpk').find("option:selected").attr("selected", false);
        $('.selectpicker').selectpicker('refresh'); 
       
      
        /*
            {
                var selectpicker_name = "";
                for(var i=0;i<$(".selectpicker").find("option:selected").length;i++){
                    if(i==0){
                        selectpicker_name = $(".selectpicker").find("option:selected")[i].innerText;//是这个，innerText
                        
                    }else{
                        selectpicker_name = selectpicker_name + ',' + $(".selectpicker").find("option:selected")[i].innerText;
        
                    }
                }//试试
                $(".selectpicker").selectpicker('deselectAll');
            }
            var message1 = $("#inputMessage1").val();
            var message2 = $("#inputMessage2").val();

            selectpicker_name =  message1 + selectpicker_name + message2;
            socket.send(selectpicker_name);
    
            console.log(selectpicker_name);
            console.log($("#inputMessage1").val());
            console.log($("#inputMessage2").val());
            */
    }
        

   
});
