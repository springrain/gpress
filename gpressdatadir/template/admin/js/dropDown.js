// 下拉菜单选封装
function DropDown(settings){
	this.settings = settings
	this.init()
}

DropDown.prototype = {
	init: function() {
		this.create()
	},
	create: function(){
		$('.dropdown-box').remove()
		let _template = `
			<div class="dropdown-box">
				<div class="dropdown-val button">
					<input class="input-c input" type="text" value="${this.settings.default}" readonly
						placeholder="请选择" id="categoryName">
					<input type="hidden" id="currval" value="${this.settings.id}" />
					<i class="iconfont icon-tubiao-"></i>
				</div>
				<ul class="dropdown-list" id="dropdown-list">
					<!-- 下拉列表 -->
				</ul>
			</div>
		`
		$(this.settings.container).append(_template)
		this.refreshDom(this.settings)
	},
	refreshDom: function(){
		var _this = this
		$('.dropdown-val').click(function() {
			let list = _this.settings.data
			let liHtml = ''
			let curr = $('#currval').val()
			$('#dropdown-list li').remove()
			$.each(list, function(i, item) {
				if (curr == item.id) {
					liHtml +=
						`<li id="${item.id}" name="${item.menuName}" class="active">${item.menuName}</li>`
				} else {
					liHtml += `<li id="${item.id}" name="${item.menuName}">${item.menuName}</li>`
				}
			})
			$('#dropdown-list').append(liHtml)
			$('#dropdown-list').toggle()
		})
		
		$('#dropdown-list').on('click', 'li', function() {
			let curr = $(this).attr('id')
			_this.settings.default = $(this).attr('name')
			_this.settings.id = curr
			_this.settings.callBack({ 'id':curr,'name':$(this).attr('name') })
			$('#dropdown-list').hide()
			_this.create()
		})
	},
}
var dropDownInit = function(opts) {
	return new DropDown(opts)
}
window.dropDownInit = $.dropDownInit = dropDownInit