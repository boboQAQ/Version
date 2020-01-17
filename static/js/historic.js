$(function () {

    var Data
    //创建websocket链接
    var socket = new WebSocket("ws://127.0.0.1:8080/WS4");

    socket.onopen = function () {
        console.log("websocket open");
        connected = true;
    };

    socket.onclose = function () {
        console.log("websocket close");
        connected = false;
    };

    //解决iframe为子窗口刷新后跳转页面回转index默认页面的问题
    // function loadIframe(url) {
    //     //获取url链接
    //     var u = window.location.href;
    //     //因为每次获取的链接中都有之前的旧锚点，
    //     //所以需要把#之后的旧锚点去掉再来加新的锚点（即传入的url参数）
    //     var end = u.indexOf("#");
    //     var rurl = u.substring(0,end);
    //     //设置新的锚点
    //     window.location.href = rurl + "#" + url;
    // }

    // var hash = location.hash;
    // var url = hash.substring(1,hash.length);
    // $("#iframe").attr("src", url); //别忘了加id


    //将不足10的补齐0
    function getzf(num) {
        if (parseInt(num) < 10) {
            num = '0' + num;
        }
        return num;
    }
    //将此类2019-12-06T09:40:11Z 格式的时间转换为正常时间
    function format(date1) {
        var date = new Date(date1);
        date = date - 8 * 3600 * 1000;
        var d = new Date(date);
        var times = getzf(d.getFullYear()) + '-' + getzf((d.getMonth() + 1)) + '-'
            + getzf(d.getDate()) + ' ' + getzf(d.getHours()) + ':' + getzf(d.getMinutes())
            + ':' + getzf(d.getSeconds());

        return times;
    }

    socket.onmessage = function (event) {
        //解析json，之后初始化加载的更新页面
        var data = JSON.parse(event.data);
        var data2 = data;
        Data = data;

        var select2 = $("#slpk2");       //下面给服务列表动态添加服务 
        var dt = data2[0].servicelist;
        var len = data2[0].servicelist.length;
        for (var i = 0; i < len; i++) {
            var str1 = dt[i].servicename + "&" + dt[i].servicenumber;
            var str2 = dt[i].servicename + "发布版本号" + dt[i].servicenumber;
            select2.append("<option value='" + str1 + "'>" + str2 + "</option>");
        }
        select2.selectpicker('refresh');
        //向表格动态添加数据
        console.log(data2);
        var versiontab = $('#versiontable');
        for (var i = 0; i < data2.length; i++) {

            var list = data2[i].servicelist;
            var str1 = "";
            for (var j = 0; j < list.length; j++) {   //根据后台数据，在服务列表的下拉框中选出该大版本下挂钩的服务
                str1 += "&lt" + list[j].servicename + "--" + list[j].servicenumber + "&gt";
            }
            console.log(str1);

            versiontab.append('<tr class="success"> ' +
                '<td>' + i + '</td>' +
                '<td>' + data2[i].versionnumber + '</td>' +
                '<td>' + str1 + '</td>' +
                '<td>' + format(data2[i].issuetime) + '</td>' +
                '<td>' + format(data2[i].creattime) + '</td>' +
                '<td>' + data2[i].comment + '</td>' +
                '</tr>')
        }
        $(document).ready(function () {
            //日志跟踪分页逻辑
            var tbtotal = $("#versiontable").find("tr").length;
            //第一行是标题 默认显示前五条纪录
            if (tbtotal > 1) {
                var countstr = "当前第 1 - 5 条　共计 ";
                countstr += tbtotal - 1;
                $("#totaltr").val(tbtotal);
                countstr += " 条";
                countstr += " 第 1 页";
                $("#gadtable_info").text(countstr);
                $("#versiontable").find("tr").each(function (i) {
                    //默认显示五条
                    if (i > 5) {
                        $(this).hide();
                    }
                })
            }
        });

        console.log("revice:", data2);     //输出解析之后的后台文件
        var select = $("#slpk1");        //给下拉框定义别名
        var list = data2[0].servicelist //默认选择了第一个，所以这是它的服务列表
        var i = 1;
        for (var j = 0; j < data2.length; j++) {   //使用jQuery动态给下拉框添加option
            if (data2[j].versionnumber == "") data2[j].versionnumber = "无效版本" + i++;
            select.append("<option value='" + data2[j].id + "'>" + data2[j].versionnumber + "</option>");
        }
        for(var k = 0; k < list.length; k++) {
            var str
            str = list[k].servicename + "&" + list[k].servicenumber;
            $("#slpk2 option[value='"+str+"']").prop("selected","selected");       
        }
        $("#slpk2").selectpicker('refresh');
        select.selectpicker('refresh');   //刷新下拉框

    }
   

    //实现点击表头表格自动排序（含数字、字符串、日期）
    var tbody = document.querySelector('#versiontable').tBodies[0];
    var th = document.querySelector('#versiontable').tHead.rows[0].cells;
    var td = tbody.rows;
    for (var i = 0; i < th.length; i++) {
        th[i].flag = 1;
        th[i].onclick = function () {
            sort(this.getAttribute('data-type'), this.flag, this.cellIndex);
            this.flag = -this.flag;
        };
    };
    function sort(str, flag, n) {
        var arr = []; //存放DOM
        for (var i = 0; i < td.length; i++) {
            arr.push(td[i]);
        };
        //排序
        arr.sort(function (a, b) {
            return method(str, a.cells[n].innerHTML, b.cells[n].innerHTML) * flag;
        });
        //添加
        for (var i = 0; i < arr.length; i++) {
            tbody.appendChild(arr[i]);
        };
    };
    //排序方法
    function method(str, a, b) {
        switch (str) {
            case 'num': //数字排序
                return a - b;
                break;
            case 'string': //字符串排序
                return a.localeCompare(b);
                break;
            default:  //日期排序，IE8下'2012-12-12'这种格式无法设置时间，替换成'/'
                return new Date(a.split('-').join('/')).getTime() - new Date(b.split('-').join('/')).getTime();
        };
    };


    $('#slpk1').change(function () {
        var data = Data;
        var id = $("#slpk1").val();
        for (var j = 0; j < data.length; j++) {
            if (data[j].id == id) {
                var list = data[j].servicelist;
                console.log("change")
                $("#slpk2").empty();
                for (var k = 0; k < list.length; k++) {
                    var str1 = list[k].servicename + "&" + list[k].servicenumber;
                    var str2 = list[k].servicename + "发布版本号" + list[k].servicenumber;
                    $("#slpk2").append("<option value='" + str1 + "'>" + str2 + "</option>");
                }
                for(var k = 0; k < list.length; k++) {
                    var str
                    str = list[k].servicename + "&" + list[k].servicenumber;
                    $("#slpk2 option[value='"+str+"']").prop("selected","selected");       
                }
                $("#slpk2").selectpicker('refresh');
                break;
            }
        }
    });



     $("#gadtable_next").click(function () {
        var nowpage = $("#nowpage").val();
        var totalpage = 0;
        //计算共多少行
        var tbtotal = $("#totaltr").val() - 1;
        if (tbtotal > 0 && tbtotal <= 5) {
            alert("已经是最后一页");
        } else if (tbtotal > 5) {
            var page = tbtotal % 5;
            totalpage = (tbtotal % 5 == 0 ? tbtotal / 5 : Math.ceil(tbtotal / 5));
            if (parseInt(nowpage) + 1 <= totalpage) {
                $("#nowpage").val(parseInt(nowpage) + 1);
                var countstr = "当前第 ";
                countstr += parseInt(nowpage) * 5 + 1;
                countstr += "-";
                countstr += nowpage == totalpage ? tbtotal : (parseInt(nowpage) + 1) * 5
                countstr += " 条　共计 ";
                countstr += tbtotal;
                countstr += " 条";
                countstr += " 第 "
                countstr += parseInt(nowpage) + 1
                countstr += " 页";
                $("#gadtable_info").text(countstr);
                $("#versiontable").find("tr").each(function (i) {
                    //当前页最后一条
                    var nowlasttr = parseInt(nowpage) * 5 + 1;
                    if (i == 0 || (nowlasttr <= i && i <= nowlasttr + 4)) {
                        $(this).show();
                    } else {
                        $(this).hide();
                    }
                })
            } else {
                alert("已是尾页");
            }
        }
    });
    $("#gadtable_previous").click(function () {
        var nowpage = $("#nowpage").val();
        var totalpage = 0;
        //计算共多少行
        var tbtotal = $("#totaltr").val() - 1;
        if (tbtotal > 0 && tbtotal <= 5) {
            alert("已经是第一页");
        } else if (tbtotal > 5) {
            var page = tbtotal % 5;
            totalpage = (tbtotal % 5 == 0 ? tbtotal / 5 : Math.ceil(tbtotal / 5));
            if (parseInt(nowpage) - 1 > 0) {
                $("#nowpage").val(parseInt(nowpage) - 1);
                var countstr = "当前第 ";
                countstr += (parseInt(nowpage) - 1) * 5 - 5;
                countstr += "-";
                countstr += (parseInt(nowpage) - 1) * 5
                countstr += " 条　共计 ";
                countstr += tbtotal;
                countstr += " 条";
                countstr += " 第 "
                countstr += parseInt(nowpage) - 1
                countstr += " 页";
                $("#gadtable_info").text(countstr);
                $("#versiontable").find("tr").each(function (i) {
                    //当前页第一条
                    var nowlasttr = parseInt(nowpage) * 5 - 5;
                    if (i == 0 || (nowlasttr - 5 <= i && i <= nowlasttr)) {
                        $(this).show();
                    } else {
                        $(this).hide();
                    }
                })
            } else {
                alert("已是首页");
            }
        }
    });


});