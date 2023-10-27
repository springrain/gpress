document.addEventListener('DOMContentLoaded', () => {
  const $navbarBurgers = document.getElementById("menuLink");
  const $layout =  document.getElementById("layout");
  const $menu =  document.getElementById("menu");
  const $main =  document.getElementById("main");
  /**
   * url
   */
  try {
    var url = window.location.pathname;
    var menus = document.getElementsByClassName("nav-addr");
    var _secondTitle="";
    for (var i in menus) {
      if (url.substring(0,url.lastIndexOf("/")) == menus[i].pathname.substring(0,url.lastIndexOf("/"))) {
        menus[i].style = "background:#1677ff"
        _secondTitle=$(menus[i]).text().trim()
        break;
      }
    }
  }catch(e){

  }
  addWebNav(_secondTitle,null);
  $navbarBurgers.addEventListener('click', () => {
    // 切换菜单的可见性
    $navbarBurgers.classList.toggle('is-active');
    $layout.classList.toggle('is-active');
    $menu.classList.toggle('is-active');
    $main.classList.toggle('is-active');
    if($main.classList.contains('is-active')){
      // console.log('window.innerWidth',window.innerWidth)
      $layout.style.minWidth = (window.innerWidth + 200) + "px";
      // console.log('$layout.style.minWidth',$layout.style.minWidth)
    }else{
      $layout.style.minWidth = (screen.width) + "px";
    }
  });
  function addWebNav(_title,_thirdTitle){
    var _template='<div class="container is-fullhd webnav"><div class="notification">' +
        '<nav class="breadcrumb" aria-label="breadcrumbs">' +
        '<ul>' +
        '<li>' +
        '<a href="#">' +
        '<span class="icon is-small">' +
        '  <i class="iconfont icon-shouye" aria-hidden="true"></i>' +
        '</span>' +
        '<span>首页</span>' +
        '</a>' +
        '</li>' +
        '<li>' +
        '<a href="#">' +
        '<span>'+_secondTitle+'</span>' +
        '</a>' +
        '</li>' +
        '</ul>' +
        '</nav>' +
        '<div class="personal">' +
        '<div class="dropdown is-hoverable">' +
        '<div class="dropdown-trigger">' +
        '<button class="button is-small" aria-haspopup="true" aria-controls="dropdown-menu4">' +
        '<span class="icon is-small">' +
        '        <i class="iconfont icon-yonghuxinxi" aria-hidden="true"></i>' +
        '      </span>' +
        '<span>个人中心</span>' +
        '</button>' +
        '</div>' +
        '<div class="dropdown-menu  is-small" id="dropdown-menu3" role="menu">' +
        '<div class="dropdown-content">' +
        '<a href="#" class="dropdown-item">' +
        '个人信息' +
        '</a>' +
        '<a href="#" class="dropdown-item">' +
        '修改密码' +
        '</a>' +
        '<hr class="dropdown-divider">' +
        '<a href="#" class="dropdown-item quit">' +
        '<span class="icon is-small">' +
        '        <i class="iconfont icon-tuichu" aria-hidden="true"></i>' +
        '      </span>' +
        '退出' +
        '</a>' +
        '</div>' +
        '</div>' +
        '</div>' +
        '</div>' +
        '' +
        '</div></div>';
    $(".section").before(_template)
  }

  /*
  // 增加form提交
  const form = document.querySelector("#gpress-form");
   
  // 如果form存在,拦截提交方法
  if (!!form) {
      form.addEventListener("submit", async (event) => {
        event.preventDefault();
        const formData = new FormData(form);
        const jsonObject = {};
        for (const [key, value] of formData.entries()) {
          jsonObject[key] = value;
        }
        const response = await fetch(form.action, {
          method: form.method,
          body: JSON.stringify(jsonObject),
          headers: {
            "Content-Type": "application/json"
          }
        });
        console.log(await response.json());
      });

    }
    */

});
