AOS.init({ duration: 800, once: true });

    document.addEventListener("DOMContentLoaded", async () => {
        //lead cites from backend 
        const citySelect = document.getElementById('heroCitySelect');
        
        try {
            if (citySelect && typeof geoApi !== "undefined") {
                const cities = await geoApi.cities();
                citySelect.innerHTML = '<option value="">اختر المدينة...</option>'; 
                
                cities.forEach(c => {
                    const opt = document.createElement('option');
                    const name = c.ville || c.name || c.Ville; 
                    opt.textContent = name;
                    opt.value = name;
                    citySelect.appendChild(opt);
                });
            }
        } catch (error) {
            console.warn("Using Fallback Cities due to API Error:", error);
            const fallback = ["Casablanca", "Rabat", "Marrakech", "Fes", "Tangier", "Agadir"];
            fallback.forEach(city => {
                const opt = document.createElement('option');
                opt.textContent = city;
                opt.value = city;
                citySelect.appendChild(opt);
            });
        }
        //lead banks walkin mazal maghdi nsta3malhouum 
        const banksContainer = document.getElementById('banks-grid');
        
        try {
            if (banksContainer) {
                const banks = await banksApi.getAll("?pageSize=6");
                 if(banks.length > 0) banksContainer.innerHTML = "";

                banks.forEach(bank => {
                    
                    let totalStock = 0;
                    if(bank.bloodGroups) {
                        bank.bloodGroups.forEach(bg => totalStock += bg.quantity);
                    }

                    let statusClass = "bg-green-100 text-green-700";
                    let statusText = "متوفر";

                    if(totalStock < 10) {
                        statusClass = "bg-red-100 text-red-700 animate-pulse";
                        statusText = "حرج";
                    } else if(totalStock < 50) {
                        statusClass = "bg-orange-100 text-orange-700";
                        statusText = "منخفض";
                    }

                    const card = `
                        <div class="bg-white border border-slate-100 rounded-2xl p-6 hover:shadow-lg transition">
                            <div class="flex justify-between items-start mb-4">
                                <span class="px-3 py-1 rounded-full text-xs font-bold ${statusClass}">${statusText}</span>
                                <i class="fa-solid fa-hospital text-slate-300"></i>
                            </div>
                            <h3 class="font-bold text-lg text-slate-900">${bank.name}</h3>
                            <p class="text-sm text-slate-500 mb-4"><i class="fa-solid fa-location-dot ml-1 text-morocco-red"></i> ${bank.city}</p>
                            
                            <div class="border-t border-slate-50 pt-3 flex justify-between text-xs text-slate-400">
                                <span>آخر تحديث:</span>
                                <span>${new Date(bank.lastUpdated).toLocaleDateString('ar-MA')}</span>
                            </div>
                        </div>`;
                    
                    banksContainer.innerHTML += card;
                });
            }
        } catch (error) {
            console.error("Using Fallback Banks:", error);
            banksContainer.innerHTML = `
                <div class="bg-white border border-slate-100 rounded-2xl p-6 text-center text-slate-400 col-span-full">
                    <i class="fa-solid fa-server mb-2 text-2xl"></i>
                    <p>لا يمكن الاتصال بالخادم حالياً.</p>
                </div>
            `;
        }
    });