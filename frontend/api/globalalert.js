 function showAutoAlert(message, type = 'error') {
    const div = document.createElement('div');
    console.log("dkhaal hna ghir mabghaach ")
    // Design
    div.style.position = 'fixed';
    div.style.top = '20px';
    div.style.left = '50%';
    div.style.transform = 'translateX(-50%)';
    div.style.padding = '15px 25px';
    // Hna fin kaytbeddel loun (Khder ola 7mer)
    div.style.backgroundColor = type === 'success' ? '#28a745' : '#dc3545'; 
    div.style.color = 'white';
    div.style.borderRadius = '8px';
    div.style.boxShadow = '0 4px 12px rgba(0,0,0,0.15)';
    div.style.zIndex = '9999';
    div.style.fontSize = '16px';
    div.style.fontWeight = '500';
    div.style.opacity = '0';
    div.style.transition = 'opacity 0.5s ease-in-out';
    
    div.innerText = message;
    
    document.body.appendChild(div);

    // Tban
    requestAnimationFrame(() => {
        div.style.opacity = '1';
    });

    // Tghber men be3d 3 seconds
    setTimeout(() => {
        div.style.opacity = '0';
        setTimeout(() => div.remove(), 500);
    }, 3000);
}
window.showAutoAlert = showAutoAlert; // assign the function itself