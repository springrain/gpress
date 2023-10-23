function initPage(){
	return $.pageInit({
			container: '.page', //容器：默认page
			countPage: $('#pageCount').val(), //一共有多少页
			showPageCount: $('#pageCount').val(), //显示多少个分页按钮
			nowPage: $('#pageNo').val(), //当前是第几页
			count: $('#totalCount').val(), //数据总数
			callBack: function (data) {
				//当前的页码
				$('#pageNo').val(data)
				$('#listForm').submit()
				//console.log('当前的页码为：' + data)
			}
		})
}