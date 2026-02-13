document.addEventListener("DOMContentLoaded", async () => {
    
    // ----------------------------------------------------
    //  INIT MAP (Leaflet)
    // ----------------------------------------------------
    // Center map on Morocco initially
    const map = L.map('map', { zoomControl: false }).setView([31.7917, -7.0926], 6);
    
    // Add CartoDB Layer (Clean Map Style)
    L.tileLayer('https://{s}.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}{r}.png', {
        attribution: '&copy; OpenStreetMap contributors &copy; CARTO',
        maxZoom: 19
    }).addTo(map);

    // ----------------------------------------------------
    //  CUSTOM ICONS (Markers Style)
    // ----------------------------------------------------
    //  User Icon (Blue Pulse)
    const userIcon = L.divIcon({
        className: 'user-marker',
        html: '<div class="w-4 h-4 bg-blue-500 rounded-full border-2 border-white shadow-lg relative"><div class="absolute inset-0 bg-blue-400 rounded-full animate-ping opacity-75"></div></div>',
        iconSize: [20, 20],
        iconAnchor: [10, 10]
    });

    // Bank Icon (Red Hospital)
    const bankIcon = L.divIcon({
        html: `<div class="w-10 h-10 bg-white rounded-full flex items-center justify-center shadow-lg border-2 border-[#C1272D] text-[#C1272D] relative transition hover:scale-110">
                 <i class="fa-solid fa-hospital text-lg"></i>
                 <span class="absolute -top-1 -right-1 flex h-3 w-3">
                    <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-red-400 opacity-75"></span>
                    <span class="relative inline-flex rounded-full h-3 w-3 bg-red-500"></span>
                 </span>
               </div>`,
        className: 'bg-transparent',
        iconSize: [40, 40],
        iconAnchor: [20, 40],
        popupAnchor: [0, -45]
    });

    try {
        // ----------------------------------------------------
        //  FETCH USER DATA (Profile)
        // ----------------------------------------------------
        const user = await accountApi.profile();
        
        const firstName = user.name.split(' ')[0] || "Hero";
        
        // Update UI Elements
        const welcomeEl = document.getElementById('welcomeName');
        const navNameEl = document.getElementById('navName');
        const navAvatar = document.getElementById('navAvatar');

        if(welcomeEl) welcomeEl.innerText = firstName;
        if(navNameEl) navNameEl.innerText = firstName;
        
        // Avatar Fallback Logic
        if(navAvatar) {
            if (user.photo_url && user.photo_url !== "") {
                navAvatar.src = user.photo_url;
            } else {
                navAvatar.src = `https://ui-avatars.com/api/?name=${user.name}&background=C1272D&color=fff`;
            }
        }

        // ----------------------------------------------------
        //  GET USER GPS LOCATION
        // ----------------------------------------------------
        if (navigator.geolocation) {
            navigator.geolocation.getCurrentPosition(
                (pos) => {
                    const { latitude, longitude } = pos.coords;
                    
                    map.setView([latitude, longitude], 14);
                    
                    L.marker([latitude, longitude], { icon: userIcon })
                     .addTo(map)
                     .bindPopup("<b>You are here</b>");
                },
                (err) => {
                    console.warn("GPS Access Denied or Error:", err.message);
                }
            );
        }

        // ----------------------------------------------------
        //  FETCH & PLOT BLOOD BANKS
        // ----------------------------------------------------
        // Using  apiRequest directly
        const banks = await apiRequest('/banks?pageSize=50', 'GET');
            
        // console.log("📍 Loaded Banks:", banks.length)

        banks.forEach(bank => {
            // console.log(`Bank: ${bank.name}, Lat: ${bank.latitude}, Lng: ${bank.longitude}`)
            if (bank.latitude && bank.longitude && bank.latitude !== 0) {
                
                let totalBags = 0;
                if (bank.bloodGroups) {
                    bank.bloodGroups.forEach(bg => totalBags += bg.quantity);
                }

                let stockBadge = totalBags < 20 
                    ? `<span class="bg-red-100 text-red-700 text-[10px] font-bold px-2 py-0.5 rounded-full animate-pulse">Critical Stock (${totalBags})</span>`
                    : `<span class="bg-green-100 text-green-700 text-[10px] font-bold px-2 py-0.5 rounded-full">Stable (${totalBags})</span>`;

                let bloodGrid = '';
                if(bank.bloodGroups) {
                    bloodGrid = bank.bloodGroups.slice(0, 4).map(bg => 
                        `<div class="bg-slate-50 rounded p-1 text-center border border-slate-100">
                            <span class="block text-[10px] font-bold text-slate-700">${bg.group}</span>
                            <span class="block text-[9px] text-slate-400 font-mono">${bg.quantity}</span>
                         </div>`
                    ).join('');
                }

                const popupContent = `
                    <div class="text-center min-w-[180px] font-sans p-1">
                        <h3 class="font-bold text-slate-900 text-sm mb-1">${bank.name}</h3>
                        
                        <p class="text-xs text-slate-500 mb-2 flex justify-center items-center gap-1">
                            <i class="fa-solid fa-location-dot text-[#C1272D]"></i> ${bank.city}
                        </p>
                        
                        <div class="mb-3 flex justify-center">${stockBadge}</div>

                        <div class="grid grid-cols-4 gap-1 mb-3">
                            ${bloodGrid}
                        </div>

                        <a href="https://www.google.com/maps/dir/?api=1&destination=${bank.latitude},${bank.longitude}" 
                           target="_blank"
                           class="block w-full bg-slate-900 text-white text-xs py-2 rounded-lg hover:bg-[#C1272D] transition font-bold flex items-center justify-center gap-2">
                           <span>Get Directions</span> 
                           <i class="fa-solid fa-arrow-up-right-from-square"></i>
                        </a>
                    </div>
                `;

                // Add Marker to Map
                L.marker([bank.latitude, bank.longitude], { icon: bankIcon })
                 .addTo(map)
                 .bindPopup(popupContent);
            }
        });

    } catch (error) {
        console.error("Dashboard Error:", error);
        
        // Handle Session Expiry
        if(error.message.includes("Unauthorized") || error.message.includes("401")) {
             localStorage.removeItem('lifeline_token');
             window.location.href = 'login.html';
        }
    }
});