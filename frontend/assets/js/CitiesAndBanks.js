// --- BACKEND LOAD SIMULATION ---
        document.addEventListener("DOMContentLoaded", () => {
            // Load Cities
            const citySelect = document.getElementById('heroCitySelect');
            if (citySelect && typeof geoApi !== "undefined") {
                geoApi.cities().then((cities) => {
                    cities.forEach(city => {
                        const opt = document.createElement('option');
                        opt.textContent = city.name;
                        opt.value = city.name;
                        citySelect.appendChild(opt);
                    });
                }).catch(() => {
                    const fallback = ["Casablanca", "Rabat", "Marrakech", "Fes", "Tangier", "Agadir"];
                    fallback.forEach(city => {
                        const opt = document.createElement('option');
                        opt.textContent = city;
                        opt.value = city;
                        citySelect.appendChild(opt);
                    });
                });
            }

            // Load Banks
            const banksContainer = document.getElementById('banks-grid');
            const banksData = [
                { name: "مركز الرباط", city: "الرباط", stock: "Low", type: "O+" },
                { name: "مركز ابن رشد", city: "الدار البيضاء", stock: "Critical", type: "AB-" },
                { name: "مستشفى محمد السادس", city: "مراكش", stock: "Stable", type: "A+" }
            ];

            banksContainer.innerHTML = "";
            banksData.forEach(bank => {
                let colorClass = "bg-green-100 text-green-700";
                if(bank.stock === "Critical") colorClass = "bg-red-100 text-red-700 animate-pulse";
                if(bank.stock === "Low") colorClass = "bg-orange-100 text-orange-700";

                const card = `
                    <div class="bg-white border border-slate-100 rounded-2xl p-6 hover:shadow-lg transition">
                        <div class="flex justify-between items-start mb-4">
                             <span class="px-3 py-1 rounded-full text-xs font-bold ${colorClass}">${bank.stock}</span>
                             <i class="fa-solid fa-hospital text-slate-300"></i>
                        </div>
                        <h3 class="font-bold text-lg text-slate-900">${bank.name}</h3>
                        <p class="text-sm text-slate-500">${bank.city}</p>
                    </div>`;
                banksContainer.innerHTML += card;
            });
        });
