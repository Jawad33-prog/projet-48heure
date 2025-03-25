document.addEventListener("DOMContentLoaded", function () {
    const button = document.getElementById("view-wines");
    if (button) {
      button.addEventListener("click", function () {
        window.location.href = "wines.html";
      });
    }
  });