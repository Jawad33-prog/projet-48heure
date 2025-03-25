document.addEventListener("DOMContentLoaded", function () {
    const button = document.getElementById("view-wines");
    if (button) {
        button.addEventListener("click", function () {
            window.location.href = "wines.html";
        });
    }

    const carousel = document.querySelector('.carousel');
    const images = carousel.querySelectorAll('img');

    let currentIndex = 0;

    function moveCarousel() {
        currentIndex = (currentIndex + 1) % images.length;

        carousel.scrollLeft = images[0].clientWidth * currentIndex; 
    }

    function updateCarousel() {
        images.forEach((image, index) => {
            const distance = Math.abs(carousel.scrollLeft - image.offsetLeft);

            if (distance < 150) { 
                image.style.filter = "brightness(100%)";
            } else {
                image.style.filter = "brightness(60%)"; 
            }
        });
    }

    carousel.addEventListener('scroll', updateCarousel);

    setInterval(moveCarousel, 3000);

    updateCarousel();
});
