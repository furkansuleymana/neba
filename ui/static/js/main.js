document.addEventListener("DOMContentLoaded", function () {
  // Highlight current page in navigation
  const currentPath = window.location.pathname;
  const navLinks = document.querySelectorAll(".main-nav a");

  navLinks.forEach((link) => {
    if (link.getAttribute("href") === currentPath) {
      link.classList.add("active");
    }
  });

  // Add fadeIn effect to the main content
  const mainContent = document.querySelector("main");
  if (mainContent) {
    mainContent.style.opacity = "0";
    mainContent.style.transition = "opacity 0.5s ease-in-out";

    setTimeout(() => {
      mainContent.style.opacity = "1";
    }, 100);
  }

  // Add click handlers for save buttons on the discover devices page
  const saveButtons = document.querySelectorAll("button[data-serial]");
  saveButtons.forEach((button) => {
    button.addEventListener("click", function () {
      alert(`This feature is not yet implemented.`);
      // TODO: Implement AJAX save functionality when backend is ready
    });
  });

  // Log that app is initialized
  console.log("Neba application initialized");
});
