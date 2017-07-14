;(function($) {
    (function mLongPolling() {
        setTimeout(function() {
            var container;
            if (document.getElementById("mlist")) {
                container = $("#mlist");
            } else {
                return true;
            };
            $.ajax({
                url: "/api/messages/",
                success: function(data) {
                    var s = "";
                    if (data) {
                        $.each(data, function(idx, item) {
                            var avatar = "";
                            if (item.Avatar) {
                                avatar = item.Avatar;
                            } else {
                                avatar = '/identicon/' + item.Sender + '/100/default.png';
                            }
                            s = s + '<li class="list-group-item"><div class="clearfix"><a href="/user/' + item.Sender + '/" class="pull-left thumb-sm avatar b-3x m-r"><img src="' + avatar + '" class="img-circle"></a><small class="pull-right">' + jsDateDiff(item.Created) + '</small><div class="clear"><div class="h3 m-t-xs m-b-xs"><a href="/connect/' + item.Sender + '/">' + item.Sender + '</a> <i class="fa fa-circle text-success pull-right text-xs m-t-sm"></i></div><small class="text-muted">' + item.Content + '</small></div></div></li>';
                        });
                    } else {
                        s = '<li class="list-group-item"><div class="clearfix"><div class="clear"><div class="h3 m-t-xs m-b-xs">尚无消息</div><small class="text-muted">你尚未收到好友消息..</small></div></div></li>';
                    };
                    container.html(s);
                },
                dataType: "json",
                complete: mLongPolling
            });
        }, 30000);
    })();
})(jQuery);