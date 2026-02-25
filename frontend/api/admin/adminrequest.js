document.addEventListener("DOMContentLoaded", loadAdminRequests);
// Load Requests
async function loadAdminRequests() {
    const tbody = document.getElementById('reqTable');
    tbody.innerHTML = '<tr><td colspan="5" class="p-4 text-center"><i class="fa-solid fa-spinner fa-spin"></i></td></tr>';

    try {
        const requests = await apiRequest('/admin/requests', 'GET');
        // console.log(requests)
        tbody.innerHTML = '';

        requests.forEach(req => {
            let statusColor = "bg-orange-100 text-orange-700";
            if (req.status === "Approved") statusColor = "bg-green-100 text-green-700";
            if (req.status === "Rejected") statusColor = "bg-red-100 text-red-700";

            let actions = '';
            if (req.status === 'Pending') {
                actions = `
                    <button onclick="updateStatus(${req.id}, 'Approved')" class="text-green-600 hover:bg-green-50 px-3 py-1 rounded font-bold transition">قبول</button>
                    <button onclick="updateStatus(${req.id}, 'Rejected')" class="text-red-600 hover:bg-red-50 px-3 py-1 rounded font-bold transition">رفض</button>
                `;
            } else {
                actions = `<span class="text-slate-400 text-xs">تمت المعالجة</span>`;
            }

            const row = `
                <tr class="hover:bg-slate-50 transition">
                    <td class="p-4 font-bold text-slate-900">${req.user?.name || 'Unknown'}</td>
                    <td class="p-4"><span class="font-black text-slate-800 bg-slate-100 px-2 rounded">${req.blood_type}</span></td>
                    <td class="p-4 text-slate-500">${req.hospital_name}</td>
                    <td class="p-4"><span class="px-2 py-1 rounded text-xs font-bold ${statusColor}">${req.status}</span></td>
                    <td class="p-4 flex gap-2">${actions}</td>
                </tr>
            `;
            tbody.innerHTML += row;
        });
    } catch (e) {
        console.error(e,"there is many problems ");
     window.showAutoAlert("Error: " + e,"error");
    }
}

// Action Function
 async function updateStatus(id, newStatus) {
    if(!confirm(`هل أنت متأكد من تغيير الحالة إلى ${newStatus}؟`)) return;
    
    try {
        await apiRequest(`/admin/requests/${id}`, 'PUT', { status: newStatus });
        loadAdminRequests(); 
        // alert("تم تحديث الحالة بنجاح");
    window.showAutoAlert("تم تحديث الحالة بنجاح","success");
    } catch (e) {
        // alert("خطأ: " + e.message);
    window.showAutoAlert("there error"+e.message,"error");
    }
}

// Init