document.addEventListener("DOMContentLoaded", function () {
  // Add click handlers for save buttons on the discover devices page
  const saveButtons = document.querySelectorAll("button[data-serial]");
  saveButtons.forEach((button) => {
    button.addEventListener("click", function () {
      // Create a simple modal dialog without Bootstrap
      const modalContainer = document.createElement("dialog");
      modalContainer.id = "notImplementedModal";

      modalContainer.innerHTML = `
        <article>
          <header>
            <h3>Not Implemented</h3>
          </header>
          <p>This feature is not yet implemented.</p>
          <footer>
            <button id="closeModalBtn">Close</button>
          </footer>
        </article>
      `;

      // Add modal to the DOM
      document.body.appendChild(modalContainer);

      // Show the modal (using the dialog element's built-in show() method)
      modalContainer.showModal();

      // Close on button click
      document
        .getElementById("closeModalBtn")
        .addEventListener("click", function () {
          modalContainer.close();
          document.body.removeChild(modalContainer);
        });

      // Also close on click outside the modal
      modalContainer.addEventListener("click", function (event) {
        if (event.target === modalContainer) {
          modalContainer.close();
          document.body.removeChild(modalContainer);
        }
      });
    });
  });

  // Log that app is initialized
  console.log("Neba application initialized with Pico CSS");
});
