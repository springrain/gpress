document.addEventListener('DOMContentLoaded', () => {
  const $navbarBurgers = document.getElementById("menuLink");
  const $layout =  document.getElementById("layout");
  const $menu =  document.getElementById("menu");
  const $main =  document.getElementById("main");

  $navbarBurgers.addEventListener('click', () => {
    // 切换菜单的可见性
    $navbarBurgers.classList.toggle('is-active');
    $layout.classList.toggle('is-active');
    $menu.classList.toggle('is-active');
    $main.classList.toggle('is-active');
    if($main.classList.contains('is-active')){
      // console.log('window.innerWidth',window.innerWidth)
      $layout.style.minWidth = (window.innerWidth + 160) + "px";
      // console.log('$layout.style.minWidth',$layout.style.minWidth)
    }else{
      $layout.style.minWidth = (screen.width) + "px";
    }
  });
});
