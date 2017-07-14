/* https://github.com/balupton/jquery-scrollto */(function(name,context,definition){if(typeof module!="undefined"&&module.exports){module.exports=definition()}else{if(typeof define=="function"&&define.amd){define(definition)}else{context[name]=definition()}}})("jquery-scrollto",this,function(){var jQuery,$,ScrollTo;jQuery=$=window.jQuery||require("jquery");$.propHooks.scrollTop=$.propHooks.scrollLeft={get:function(elem,prop){var result=null;if(elem.tagName==="HTML"||elem.tagName==="BODY"){if(prop==="scrollLeft"){result=window.scrollX}else{if(prop==="scrollTop"){result=window.scrollY}}}if(result==null){result=elem[prop]}return result}};$.Tween.propHooks.scrollTop=$.Tween.propHooks.scrollLeft={get:function(tween){return $.propHooks.scrollTop.get(tween.elem,tween.prop)},set:function(tween){if(tween.elem.tagName==="HTML"||tween.elem.tagName==="BODY"){tween.options.bodyScrollLeft=(tween.options.bodyScrollLeft||window.scrollX);tween.options.bodyScrollTop=(tween.options.bodyScrollTop||window.scrollY);if(tween.prop==="scrollLeft"){tween.options.bodyScrollLeft=Math.round(tween.now)}else{if(tween.prop==="scrollTop"){tween.options.bodyScrollTop=Math.round(tween.now)}}window.scrollTo(tween.options.bodyScrollLeft,tween.options.bodyScrollTop)}else{if(tween.elem.nodeType&&tween.elem.parentNode){tween.elem[tween.prop]=tween.now}}}};ScrollTo={config:{duration:400,easing:"swing",callback:undefined,durationMode:"each",offsetTop:0,offsetLeft:0},configure:function(options){$.extend(ScrollTo.config,options||{});return this},scroll:function(collections,config){var collection,$container,container,$target,$inline,position,containerTagName,containerScrollTop,containerScrollLeft,containerScrollTopEnd,containerScrollLeftEnd,startOffsetTop,targetOffsetTop,targetOffsetTopAdjusted,startOffsetLeft,targetOffsetLeft,targetOffsetLeftAdjusted,scrollOptions,callback;collection=collections.pop();$container=collection.$container;$target=collection.$target;containerTagName=$container.prop("tagName");$inline=$("<span/>").css({"position":"absolute","top":"0px","left":"0px"});position=$container.css("position");$container.css({position:"relative"});$inline.appendTo($container);startOffsetTop=$inline.offset().top;targetOffsetTop=$target.offset().top;targetOffsetTopAdjusted=targetOffsetTop-startOffsetTop-parseInt(config.offsetTop,10);startOffsetLeft=$inline.offset().left;targetOffsetLeft=$target.offset().left;targetOffsetLeftAdjusted=targetOffsetLeft-startOffsetLeft-parseInt(config.offsetLeft,10);containerScrollTop=$container.prop("scrollTop");containerScrollLeft=$container.prop("scrollLeft");$inline.remove();$container.css({position:position});scrollOptions={};callback=function(event){if(collections.length===0){if(typeof config.callback==="function"){config.callback()}}else{ScrollTo.scroll(collections,config)}return true};if(config.onlyIfOutside){containerScrollTopEnd=containerScrollTop+$container.height();containerScrollLeftEnd=containerScrollLeft+$container.width();if(containerScrollTop<targetOffsetTopAdjusted&&targetOffsetTopAdjusted<containerScrollTopEnd){targetOffsetTopAdjusted=containerScrollTop}if(containerScrollLeft<targetOffsetLeftAdjusted&&targetOffsetLeftAdjusted<containerScrollLeftEnd){targetOffsetLeftAdjusted=containerScrollLeft}}if(targetOffsetTopAdjusted!==containerScrollTop){scrollOptions.scrollTop=targetOffsetTopAdjusted}if(targetOffsetLeftAdjusted!==containerScrollLeft){scrollOptions.scrollLeft=targetOffsetLeftAdjusted}if($container.prop("scrollHeight")===$container.width()){delete scrollOptions.scrollTop}if($container.prop("scrollWidth")===$container.width()){delete scrollOptions.scrollLeft}if(scrollOptions.scrollTop!=null||scrollOptions.scrollLeft!=null){$container.animate(scrollOptions,{duration:config.duration,easing:config.easing,complete:callback})}else{callback()}return true},fn:function(options){var collections,config,$container,container;collections=[];var $target=$(this);if($target.length===0){return this}config=$.extend({},ScrollTo.config,options);$container=$target.parent();container=$container.get(0);while(($container.length===1)&&(container!==document.body)&&(container!==document)){var containerScrollTop,containerScrollLeft;containerScrollTop=$container.css("overflow-y")!=="visible"&&container.scrollHeight!==container.clientHeight;containerScrollLeft=$container.css("overflow-x")!=="visible"&&container.scrollWidth!==container.clientWidth;if(containerScrollTop||containerScrollLeft){collections.push({"$container":$container,"$target":$target});$target=$container}$container=$container.parent();container=$container.get(0)}collections.push({"$container":$("html"),"$target":$target});if(config.durationMode==="all"){config.duration/=collections.length}ScrollTo.scroll(collections,config);return this}};$.ScrollTo=$.ScrollTo||ScrollTo;$.fn.ScrollTo=$.fn.ScrollTo||ScrollTo.fn;return ScrollTo});
;(function($) {
    // Avoid embed thie site in an iframe of other WebSite
    if (top.location != location) {
        top.location.href = location.href;
    }

    (function(){
        // extend jQuery ajax, set xsrf token value
        var ajax = $.ajax;
        $.extend({
            ajax: function(url, options) {
                if (typeof url === 'object') {
                    options = url;
                    url = undefined;
                }
                options = options || {};
                url = options.url;
                var xsrftoken = $('meta[name=_xsrf]').attr('content');
                var headers = options.headers || {};
                var domain = document.domain.replace(/\./ig, '\\.');
                if (!/^(http:|https:).*/.test(url) || eval('/^(http:|https:)\\/\\/(.+\\.)*' + domain + '.*/').test(url)) {
                    headers = $.extend(headers, {'X-Xsrftoken':xsrftoken});
                }
                options.headers = headers;
                var callback = options.success;
                options.success = function(data){
                    if(callback){
                        callback.apply(this, arguments);
                    }
                };
                return ajax(url, options);
            }
        });

        $(document).ready(
            function() {
                $(document).scroll(function() {
                    if ($(this).scrollTop() > 720) {
                        $('#backtop').removeClass('hidden');
                    } else {
                        $('#backtop').addClass('hidden');
                    }
                });
                $('#backtop').click(function() {
                    $('body, html').animate({
                        scrollTop: 0
                    });
                    return false;
                })
                window.console && window.console.info('Copyright 2014～2016 YouGam.com');
            }
        );
    })();

    $(function() {
        $(document).on("click", "[rel=comment-reply]", function() {
            var $e = $(this).parents(".comment:first"),
                user = $e.data("user"),
                floor = $e.data("floor"),
                content = $(".reply-body-" + floor).html(),
                v = "<blockquote>" + content + "</blockquote>" + "<p>#" + floor + " @" + user + "，</p>";
            $("#post-reply").ScrollTo();
            tinyMCE.activeEditor.setContent(v + tinyMCE.activeEditor.getContent())
        })
        var $comments = $('.post-comments');

        $(window).on('hashchange', function() {
            if (/#reply\d+/.test(window.location.hash)) {
                $comments.find('.comment').removeClass('highlight');
                var $e = $(window.location.hash);
                $e.addClass('highlight');
            }
        });
        $(window).trigger('hashchange');

        $(document).on("click", "[rel=quick-reply]", function() {
            var $e = $(this),
                target = $e.data("target"),
                user = $e.data("user"),
                v = "@" + user + "，",
                tar = $("#"+target);
            tar.ScrollTo();
            tar.removeClass("hide");
            var c =$("input[name='comment']");
            c.val(v);
        })
    });
})(jQuery);

//读取cookies
function getCookie(name) {
    var arr, reg = new RegExp("(^| )" + name + "=([^;]*)(;|$)");
    if (arr = document.cookie.match(reg))
        return unescape(arr[2]);
    else
        return null;
};

function showholder(id) {
    var s = document.getElementById(id);
    if ((s.value == -1) || (s.value == -2) || (s.value == -3)){
        document.getElementById("excerptor").style.display = "block";
    } else {
        document.getElementById("excerptor").style.display = "none";
    };
};

function jsDateDiff(publishTime) {
    var d_minutes, d_hours, d_days;
    var timeNow = parseInt(new Date().getTime() / 1000);
    var d;
    d = timeNow - publishTime;
    d_days = parseInt(d / 86400);
    d_hours = parseInt(d / 3600);
    d_minutes = parseInt(d / 60);

    if (d_days > 0 && d_days < 4) {
        return d_days + " 天之前";
    } else if (d_days <= 0 && d_hours > 0) {
        return d_hours + " 小时之前";
    } else if (d_hours <= 0 && d_minutes > 0) {
        return d_minutes + " 分钟之前";
    } else {
        var s = new Date(publishTime * 1000);
        // s.getFullYear()+"年";
        return (s.getMonth() + 1) + "月" + s.getDate() + "日";
    }
};

function getUnixTime(dateStr) {
    var newstr = dateStr.replace(/-/g, '/');
    var date = new Date(newstr);
    var time_str = date.getTime().toString();
    return time_str.substr(0, 10);
};