const refresh = document.querySelector("#refresh");
const notice = document.querySelector("#notice");

refresh.addEventListener("click", () => {
  notice.textContent = `最近采样：${new Date().toLocaleTimeString("zh-CN")} 更新`;
});
