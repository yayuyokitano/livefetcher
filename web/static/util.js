const userButton = document.querySelectorAll(".user-button");
const userPopup = document.querySelectorAll(".user-popup");

userButton.forEach((btn) => {
	btn.addEventListener("click", () => {
		btn.nextElementSibling.show();
	});
});

document.addEventListener("click", (e) => {
	if (e.target.closest(".dialog-button-wrapper")) {
		return;
	}

	userPopup.forEach((popup) => {
		popup.close();
	});
});