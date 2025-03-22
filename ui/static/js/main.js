document.addEventListener("DOMContentLoaded", function () {
  // Save button handlers
  $("button[data-serial]").on("click", function () {
    $(".ui.modal").modal("show");
  });

  // Modal template
  const modalHTML = `
    <div class="ui modal">
      <div class="header">Not Implemented</div>
      <div class="content">
        <p>This feature is not yet implemented.</p>
      </div>
      <div class="actions">
        <div class="ui approve button">Close</div>
      </div>
    </div>
  `;

  // Add modal to DOM
  $("body").append(modalHTML);

  // Initialize Semantic UI components
  $(".ui.dropdown").dropdown();

  console.log("Neba application initialized with Semantic UI");
});
