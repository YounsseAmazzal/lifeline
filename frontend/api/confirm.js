export function showConfirm(message, options = {}) {
  return new Promise((resolve) => {
    const modal = document.getElementById("globalConfirmModal");
    const messageBox = document.getElementById("confirmMessage");
    const okBtn = document.getElementById("confirmOk");
    const cancelBtn = document.getElementById("confirmCancel");

    messageBox.textContent = message;

    // Optional customization
    okBtn.textContent = options.okText || "تأكيد";
    cancelBtn.textContent = options.cancelText || "إلغاء";

    modal.classList.add("active");

    const close = (result) => {
      modal.classList.remove("active");
      okBtn.onclick = null;
      cancelBtn.onclick = null;
      resolve(result);
    };

    okBtn.onclick = () => close(true);
    cancelBtn.onclick = () => close(false);

    // Close when clicking outside
    modal.onclick = (e) => {
      if (e.target === modal) close(false);
    };
  });
}