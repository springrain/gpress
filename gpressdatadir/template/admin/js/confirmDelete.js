function deleteFunc(id,url) {
		$.confirm({
			useBootstrap: false,
			title: '提示',
			content: '确认删除？',
			boxWidth: '400px',
			type: 'red',
			buttons: {
				'确定': function () {
					$.ajax({
						type: 'Post',
						url: url, //'{{basePath}}admin/{{$path}}/delete'
						data: { id: id },
						success: function (res) {
							if (res.statusCode === 1) {
								location.reload()
							}
						}
					})
				},
				'取消': function () { }
			}
		})
	}