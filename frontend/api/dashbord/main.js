document.addEventListener("DOMContentLoaded", async () => {
    
    const map = L.map('map', { 
        zoomControl: false, 
        attributionControl: false 
    }).setView([31.7917, -7.0926], 6);

    L.tileLayer('https://{s}.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}{r}.png', {
        maxZoom: 19
    }).addTo(map);

    L.control.zoom({ position: 'bottomright' }).addTo(map);

    const userIcon = L.divIcon({
        className: 'bg-transparent',
        html: `<div class="marker-pulse"></div>`,
        iconSize: [20, 20],
        iconAnchor: [10, 10]
    });

    const bankIcon = L.divIcon({
        className: 'bg-transparent',
        html: `<div class="w-10 h-10 bg-white rounded-full shadow-xl border-[3px] border-[#C1272D] flex items-center justify-center text-[#C1272D] hover:scale-110 transition-transform duration-300">
                 <i class="fa-solid fa-hospital text-sm"></i>
               </div>`,
        iconSize: [40, 40],
        iconAnchor: [20, 40],
        popupAnchor: [0, -45]
    });

    try {
        const [user, banks] = await Promise.all([
            accountApi.profile(),
            apiRequest('/banks?pageSize=50', 'GET')
        ]);

        const firstName = user.name.split(' ')[0] || "Hero";
        document.getElementById('navName').innerText = firstName;
        if(user.photo_url) document.getElementById('navAvatar').src = user.photo_url;

        if (navigator.geolocation) {
            navigator.geolocation.getCurrentPosition(pos => {
                const { latitude, longitude } = pos.coords;
                map.setView([latitude, longitude], 13);
                L.marker([latitude, longitude], { icon: userIcon }).addTo(map);
            });
        }

        const listContainer = document.getElementById('banksList');
        listContainer.innerHTML = ''; 
            if (banks.length === 0) {
        listContainer.innerHTML = `
            <div class="text-center text-slate-400 py-8">
                <i class="fa-solid fa-hospital-user text-2xl mb-2"></i>
                <p class="text-xs">No nearby banks found.</p>
            </div>`;
    }
        banks.forEach(bank => {
            if (bank.latitude && bank.longitude) {
                const marker = L.marker([bank.latitude, bank.longitude], { icon: bankIcon })
                    .addTo(map)
                    .bindPopup(createPopup(bank));

                const listItem = createListItem(bank);
                listItem.addEventListener('click', () => {
                    map.flyTo([bank.latitude, bank.longitude], 15, { duration: 1.5 });
                    marker.openPopup();
                });
                listContainer.appendChild(listItem);
            }
        });

    } catch (error) {
        console.error("Dashboard Error:", error);
        if(error.message.includes("Unauthorized")) window.location.href = 'login.html';
    }
});


function createPopup(bank) {
    let stockTotal = 0;
    bank.bloodGroups?.forEach(bg => stockTotal += bg.quantity);
    
    const statusColor = stockTotal < 20 ? 'text-red-600' : 'text-green-600';
    const statusText = stockTotal < 20 ? 'Critical' : 'Stable';

    return `
        <div class="text-center font-sans p-1">
            <h3 class="font-bold text-slate-900 text-sm mb-1">${bank.name}</h3>
            <p class="text-[10px] text-slate-500 mb-2">${bank.city}</p>
            <div class="bg-slate-50 rounded-lg p-2 mb-2 border border-slate-100">
                <span class="text-xs font-bold ${statusColor}">${statusText} Stock</span>
                <span class="block text-[10px] text-slate-400">${stockTotal} bags available</span>
            </div>
            <a href="https://maps.google.com/?q=${bank.latitude},${bank.longitude}" target="_blank" class="block w-full bg-[#C1272D] text-white text-xs py-1.5 rounded-md font-bold">Directions</a>
        </div>
    `;
}

function createListItem(bank) {
    const div = document.createElement('div');
    div.className = "bg-white p-3 rounded-xl border border-slate-100 hover:border-[#C1272D] hover:shadow-md transition cursor-pointer flex justify-between items-center group";
    
    let stockTotal = 0;
    bank.bloodGroups?.forEach(bg => stockTotal += bg.quantity);

    div.innerHTML = `
        <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-slate-50 rounded-full flex items-center justify-center text-slate-400 group-hover:bg-red-50 group-hover:text-red-500 transition">
                <i class="fa-solid fa-hospital"></i>
            </div>
            <div>
                <h4 class="font-bold text-sm text-slate-800 line-clamp-1">${bank.name}</h4>
                <p class="text-xs text-slate-500">${bank.city}</p>
            </div>
        </div>
        <div class="text-right">
            <span class="block font-bold text-xs ${stockTotal < 20 ? 'text-red-600' : 'text-green-600'}">${stockTotal < 20 ? 'Critical' : 'Stable'}</span>
            <span class="text-[10px] text-slate-400">View Map</span>
        </div>
    `;
    return div;
}



// --- TABS LOGIC ---
window.switchTab = function(tabName) {
    const banksList = document.getElementById('banksList');
    const requestsList = document.getElementById('requestsList');
    const tabBanks = document.getElementById('tab-banks');
    const tabRequests = document.getElementById('tab-requests');

    if (tabName === 'banks') {
        banksList.classList.remove('hidden');
        requestsList.classList.add('hidden');
        
        tabBanks.className = "flex-1 pb-3 text-sm font-bold text-brand-red border-b-2 border-brand-red transition";
        tabRequests.className = "flex-1 pb-3 text-sm font-bold text-slate-400 hover:text-slate-600 transition";
    } else {
        banksList.classList.add('hidden');
        requestsList.classList.remove('hidden');
        
        tabBanks.className = "flex-1 pb-3 text-sm font-bold text-slate-400 hover:text-slate-600 transition";
        tabRequests.className = "flex-1 pb-3 text-sm font-bold text-brand-red border-b-2 border-brand-red transition";
    }
}

// --- RENDER REQUESTS ---
function loadRequests() {
    const container = document.getElementById('requestsList');
    
    //ghadi njibhouum min api min 
    const requests = [
        { name: "Fatima Z.", hospital: "Ibn Sina", blood: "O+", dist: "2km", urgent: true },
        { name: "Karim M.", hospital: "Cheikh Zaid", blood: "AB-", dist: "5km", urgent: false },
        { name: "Accident #402", hospital: "Military Hosp.", blood: "A-", dist: "8km", urgent: true }
    ];

    container.innerHTML = requests.map(req => `
        <div class="bg-white p-3 rounded-xl border ${req.urgent ? 'border-red-100 bg-red-50/30' : 'border-slate-100'} flex justify-between items-center group cursor-pointer hover:shadow-md transition">
            <div class="flex items-center gap-3">
                <div class="w-10 h-10 ${req.urgent ? 'bg-red-100 text-red-600' : 'bg-slate-100 text-slate-500'} rounded-full flex items-center justify-center font-black text-xs border-2 border-white shadow-sm">
                    ${req.blood}
                </div>
                <div>
                    <h4 class="font-bold text-sm text-slate-800">${req.name}</h4>
                    <p class="text-xs text-slate-500 flex items-center gap-1">
                        <i class="fa-solid fa-hospital"></i> ${req.hospital}
                    </p>
                </div>
            </div>
            <div class="text-right">
                <span class="block font-bold text-xs text-slate-700">${req.dist}</span>
                <button class="text-[10px] bg-slate-900 text-white px-2 py-1 rounded mt-1 hover:bg-brand-red transition">Accept</button>
            </div>
        </div>
    `).join('');
}

loadRequests();

function toggleRequestModal() {
    const modal = document.getElementById('requestModal');
    modal.classList.toggle('hidden');
}
document.getElementById('requestForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const btn = e.target.querySelector('button[type="submit"]'); 
    const originalText = btn.innerHTML;
    btn.innerHTML = '<i class="fa-solid fa-circle-notch fa-spin"></i> جاري الإرسال...';
    btn.disabled = true;

    try {
        const bloodTypeInput = document.querySelector('input[name="bloodType"]:checked');
        const hospitalSelect = document.getElementById('reqHospital');
        const isUrgent = document.getElementById('reqUrgent').checked;

        if (!bloodTypeInput) {
            // alert("المرجو اختيار فصيلة الدم");
            window.showAutoAlert("المرجو اختيار مستشفى صحيح من القائمة","error");
            throw new Error("Blood type missing");
        }

        const selectedOption = hospitalSelect.options[hospitalSelect.selectedIndex];
        if (!selectedOption || selectedOption.disabled || !selectedOption.hasAttribute('data-lat')) {
            // alert("المرجو اختيار مستشفى صحيح من القائمة");
            window.showAutoAlert("المرجو اختيار مستشفى صحيح من القائمة","error");
            throw new Error("Invalid hospital selection");
        }

        const lat = parseFloat(selectedOption.getAttribute('data-lat'));
        const lng = parseFloat(selectedOption.getAttribute('data-lng'));
        const hospitalName = selectedOption.value;

        const payload = {
            bloodType: bloodTypeInput.value,
            hospitalName: hospitalName,
            isUrgent: isUrgent,
            latitude: lat,   
            longitude: lng   
        };

        // 5. Send to Backend
        await apiRequest('/requests', 'POST', payload);
        
        // Success Feedback
        // alert("تم إرسال طلبك بنجاح! سنقوم بإشعار المتبرعين القريبين من " + hospitalName)
    window.showAutoAlert("تم إرسال طلبك بنجاح! سنقوم بإشعار المتبرعين القريبين من " + hospitalName,"success");
        toggleRequestModal(); 
        e.target.reset();     
        
        hospitalSelect.innerHTML = '<option value="" disabled selected>اختر المدينة أولاً</option>';
        hospitalSelect.disabled = true;

    } catch (error) {
        if (error.message !== "Blood type missing" && error.message !== "Invalid hospital selection") {
            // alert("فشل إرسال الطلب: " + error.message);
            window.showAutoAlert("فشل إرسال الطلب: " + error.message,"error");
        }
    } finally {
        // Restore Button
        btn.innerHTML = originalText;
        btn.disabled = false;
    }
});

// Add this logic

//  Toggle Panel
window.toggleNotifPanel = function() {
    document.getElementById('notifPanel').classList.toggle('hidden');
}

// Fetch Notifications
async function checkNotifications() {
    try {
        const notifs = await apiRequest('/notifications', 'GET');
        const list = document.getElementById('notifList');
        const badge = document.getElementById('notifBadge');

        if (notifs.length > 0) {
            badge.classList.remove('hidden'); // Show Red Dot
            list.innerHTML = '';
            
            notifs.forEach(n => {
                const item = `
                    <div class="bg-red-50 p-3 rounded-xl border border-red-100 relative group">
                        <h4 class="font-bold text-sm text-red-700">${n.title}</h4>
                        <p class="text-xs text-slate-600 mt-1">${n.message}</p>
                        <span class="text-[10px] text-slate-400 block mt-2">${new Date(n.created_at).toLocaleTimeString()}</span>
                        ${!n.is_read ? '<span class="absolute top-2 right-2 w-2 h-2 bg-red-500 rounded-full"></span>' : ''}
                    </div>
                `;
                list.innerHTML += item;
            });
        }
    } catch (e) {
        // console.log("No notifications");
            window.showAutoAlert("No notifications","error");
    }
}

// Check every 10 seconds (Polling)
setInterval(checkNotifications, 10000);
// Check on Load
checkNotifications();

// --- RENDER REQUESTS ---
async function loadRequests() {
    const container = document.getElementById('requestsList');
    container.innerHTML = '<div class="text-center py-4"><i class="fa-solid fa-spinner fa-spin text-slate-400"></i></div>';

    try {
        // Fetch Approved Requests
        const requests = await apiRequest('/requests', 'GET'); 
        
        container.innerHTML = '';
        if (requests.length === 0) {
            container.innerHTML = '<div class="text-center text-slate-400 text-xs py-4">No active requests</div>';
            return;
        }

        requests.forEach(req => {
            const isUrgent = req.is_urgent;
            const card = `
                <div class="bg-white p-3 rounded-xl border ${isUrgent ? 'border-red-200 bg-red-50/50' : 'border-slate-100'} flex justify-between items-center group cursor-pointer hover:shadow-md transition">
                    <div class="flex items-center gap-3">
                        <div class="w-10 h-10 ${isUrgent ? 'bg-red-500 text-white' : 'bg-slate-100 text-slate-500'} rounded-full flex items-center justify-center font-black text-xs border-2 border-white shadow-sm">
                            ${req.blood_type}
                        </div>
                        <div>
                            <h4 class="font-bold text-sm text-slate-800">${req.user?.name || 'Patient'}</h4>
                            <p class="text-xs text-slate-500 flex items-center gap-1">
                                <i class="fa-solid fa-hospital"></i> ${req.hospital_name}
                            </p>
                        </div>
                    </div>
                    <div class="text-right">
                        ${isUrgent ? '<span class="text-[10px] text-red-600 font-bold animate-pulse">URGENT</span>' : ''}
<button onclick="acceptRequest(${req.id}, ${req.latitude}, ${req.longitude})" class="text-[10px] bg-slate-900 text-white px-3 py-1.5 rounded mt-1 hover:bg-[#C1272D] transition">
    Accept & Go
</button>                    </div>
                </div>
            `;
            container.innerHTML += card;
        });

    } catch (e) {
        console.error("Failed to load requests", e);
    }
}
async function acceptRequest(id, lat, lng) { 
    // if(!confirm("Are you sure? The patient will be notified.")) return;
    const confirmed = await showConfirm(
            "Are you sure? The patient will be notified."
        );
    
        if (!confirmed) return;
    

    try {
        await apiRequest(`/requests/${id}/accept`, 'PUT');
        
        const url = `https://www.google.com/maps/dir/?api=1&destination=${lat},${lng}`;
        window.open(url, '_blank'); 
        
        // alert("Thank you! Follow the map to the hospital.");
            window.showAutoAlert("Thank you! Follow the map to the hospital.","success");
        loadRequests(); 

    } catch (e) {
        // alert("Error: " + e.message);
     window.showAutoAlert("Error: " + e.message,"success");
    }
}
document.addEventListener("DOMContentLoaded", async () => {
    loadRequests(); 
const myReq = await apiRequest('/requests/active', 'GET');
if (myReq) {
    document.getElementById('trackingCard').classList.remove('hidden');
    document.getElementById('trackingMsg').innerText = "A donor is coming to " + myReq.hospital_name;
}
});

// --- DYNAMIC HOSPITAL SELECTOR ---

//  Load Cities 
async function initRequestModal() {
    try {
        const cities = await geoApi.cities();
        const citySelect = document.getElementById('reqCity');
        
        cities.forEach(c => {
            const opt = document.createElement('option');
            opt.value = c.ville; 
            opt.innerText = c.ville;
            citySelect.appendChild(opt);
        });
    } catch (e) {
        console.error("Failed to load cities for modal");
    }
}

//  Load Hospitals when City changes
window.loadHospitalsByCity = async function() {
    const city = document.getElementById('reqCity').value;
    const hospitalSelect = document.getElementById('reqHospital');
    
    hospitalSelect.innerHTML = '<option>جاري التحميل...</option>';
    hospitalSelect.disabled = true;

    try {
        // Fetch Banks in this City
        const banks = await apiRequest(`/banks?city=${city}&pageSize=50`, 'GET');
        // console.log(city)
        // console.log(banks)
        hospitalSelect.innerHTML = '<option value="" disabled selected>اختر المستشفى</option>';
        
        if (banks.length === 0) {
            hospitalSelect.innerHTML = '<option>لا توجد مراكز مسجلة هنا</option>';
            return;
        }

        banks.forEach(bank => {
            const opt = document.createElement('option');
            opt.value = bank.name;
            opt.innerText = bank.name;
            
            
            opt.setAttribute('data-lat', bank.latitude);
            opt.setAttribute('data-lng', bank.longitude);
            
            hospitalSelect.appendChild(opt);
        });

        hospitalSelect.disabled = false; // Unlock

    } catch (e) {
        hospitalSelect.innerHTML = '<option>خطأ في التحميل</option>';
    }
}
function showConfirm(message, options = {}) {
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
// Call init on load
document.addEventListener("DOMContentLoaded", initRequestModal);