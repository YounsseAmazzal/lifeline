
document.addEventListener("DOMContentLoaded", async () => {
    const map = L.map('map', { zoomControl: false }).setView([31.7917, -7.0926], 6);
    L.tileLayer('https://{s}.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}{r}.png', {
        attribution: '&copy; OpenStreetMap'
    }).addTo(map);

    try {
        const user = await accountApi.profile();
        const firstName = user.name.split(' ')[0];
        const welcomeEl = document.getElementById('welcomeName');
        const navNameEl = document.getElementById('navName');
        const navAvatar = document.getElementById('navAvatar');

        if(welcomeEl) welcomeEl.innerText = firstName;
        if(navNameEl) navNameEl.innerText = firstName;
        if(navAvatar && user.photo_url) navAvatar.src = user.photo_url;

        if (navigator.geolocation) {
            navigator.geolocation.getCurrentPosition(pos => {
                const { latitude, longitude } = pos.coords;
                map.setView([latitude, longitude], 14);
                
                const userIcon = L.divIcon({
                    className: 'user-marker', 
                    iconSize: [20, 20]
                });
                L.marker([latitude, longitude], { icon: userIcon }).addTo(map);
            });
        }
        //khaasni n9ad logic dyal banks walkin machi daba 
        const banks = await apiRequest('/banks', 'GET');
        // console.log(banks.Longitude)

    } catch (error) {
        console.error("Dashboard Error:", error);
        if(error.message.includes("Unauthorized")) {
             window.location.href = 'login.html';
        }
    }
});
