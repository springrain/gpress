// 基于jquery的分页控件
function Page(settings) {
	this.settings = settings
	this.init()
}
//默认配置
Page.prototype = {
	init: function() {
		this.create()
	},
	create: function() {
		var _template = `<nav class="pagination is-centered" role="navigation" aria-label="pagination">
		<div class="flex">
			<div class="lineheight40">共 ${this.settings.count} 条</div>
			<div><a class="pagination-link" aria-label="总页数">共 ${this.settings.countPage} 页</a></div>
			<div><a class="pagination-previous page_up">上一页</a></div>
			<ul class="pagination-list page_view_ul">
			</ul>
			<div><a class="pagination-next page_down">下一页</a></div>
			<div class="pagination-page flex">
				<input class="input page_input" type="text" placeholder="页码"> &nbsp;
				<button class="button page_btn">跳转</button>
			</div>
		</div>
    </nav>`
		$(this.settings.container).append(_template)
		this.refreshDom(this.settings)
		this.bindEvent()
	},
	bindEvent: function() {
		var _this = this
		//跳转首页
		$(this.settings.container).on('click', '.page_home', function() {
			var newpages = 1
			_this.settings.nowPage = newpages
			_this.settings.callBack(_this.settings.nowPage)
			_this.refreshDom(this.settings)
		})
		//跳转上一页
		$(this.settings.container).on('click', '.page_up', function() {
			var newpages = _this.settings.nowPage
			newpages--
			if (newpages < 1) {
				newpages = 1
				_this.settings.nowPage = newpages
			} else {
				_this.settings.nowPage = newpages
			}
			_this.settings.callBack(_this.settings.nowPage)
			_this.refreshDom(this.settings)
		})
		//跳转下一页
		$(this.settings.container).on('click', '.page_down', function() {
			var newpages = _this.settings.nowPage
			newpages++
			if (newpages > _this.settings.countPage) {
				newpages = _this.settings.countPage
				_this.settings.nowPage = newpages
			} else {
				_this.settings.nowPage = newpages
			}
			_this.settings.callBack(_this.settings.nowPage)
			_this.refreshDom(this.settings)
		})
		//跳转末页
		$(this.settings.container).on('click', '.page_trailer', function() {
			var newpages = _this.settings.countPage
			_this.settings.nowPage = newpages
			_this.settings.callBack(_this.settings.nowPage)
			_this.refreshDom(this.settings)
		})
		//Go跳转
		$(this.settings.container).on('click', '.page_btn', function() {
			var inputText = $('.page_input').val() - 0
			if (inputText < 1 || inputText > _this.settings.countPage) {
				$.confirm({
					useBootstrap:false,
					boxWidth: '400px',
					title: '提示',
					content: '请输入正确的页码！',
					type: 'red',
					buttons: { '确定': function() {} }
				})
			} else {
				_this.settings.nowPage = inputText
				_this.settings.callBack(_this.settings.nowPage)
				_this.refreshDom(this.settings)
				//                establish(objpage);
				//外部的ajax
			}
		})
	},
	refreshDom: function() {
		var _this = this
		$('.li').remove()
		var countPage = this.settings.countPage - 0
		var showPageCount = this.settings.showPageCount - 0
		var nowPage = this.settings.nowPage - 0
		var count = this.settings.count - 0
		var bian = (showPageCount) / 2
		// console.log(countPage,'总共多少页')
		// console.log(showPageCount,'显示多少个分页按钮')
		// console.log(nowPage,'当前是第几页')
		// console.log(bian,'bianbianbianbian')
		$('.all_data').html(count)
		$('.all_pages').html(countPage)
		var html = ''
		if (nowPage - bian <= 0) {
			for (var i = 1; i < showPageCount + 1; i++) {
				var index = i
				if (nowPage === index) {
					pageHtml =
						`<li index="${i}" class="li active"><a class="pagination-link is-current">${i}</a></li>`
				} else {
					pageHtml = `<li  index="${i}" class="li"><a class="pagination-link">${i}</a></li>`
				}
				html += pageHtml
			}
		} else if (nowPage - bian > 0 && nowPage + bian < countPage) {
			var num = nowPage - bian
			for (var i = num; i < (num + showPageCount); i++) {
				var index = i
				if (nowPage === index) {
					pageHtml = `<li index="${i}" class="li"><a class="pagination-link is-current">${i}</a></li>`
				} else {
					pageHtml = `<li index="${i}" class="li"><a class="pagination-link">${i}</a></li>`
				}
				html += pageHtml
			}
		} else if (nowPage + bian >= countPage) {
			var numAll = countPage - showPageCount + 1
			for (var i = numAll; i < (numAll + showPageCount); i++) {
				var index = i
				if (nowPage === index) {
					pageHtml = `<li index="${i}" class="li"><a class="pagination-link is-current">${i}</a></li>`
				} else {
					pageHtml = `<li index="${i}" class="li"><a class="pagination-link">${i}</a></li>`
				}
				html += pageHtml
			}
		}
		$('.page_view_ul').append(html)
		$('.li').click(function() {
			_this.settings.nowPage = $(this).attr('index') - 0
			_this.settings.callBack(_this.settings.nowPage)
			_this.refreshDom(this.settings)
		})
	}
}
var pageInit = function(opts) {
	return new Page(opts)
}
window.pageInit = $.pageInit = pageInit
