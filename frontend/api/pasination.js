let currentPage = 1;
const itemsPerPage = 6; 

async function loadBanks(page = 1) {
    try {
        const container = document.getElementById('banks-grid'); 
        if (!container) return; 

        container.innerHTML = '<div class="col-span-full text-center py-12"><i class="fa-solid fa-spinner fa-spin text-2xl text-morocco-green"></i></div>';

        const response = await fetch(`${API_URL}/banks?pageNumber=${page}&pageSize=${itemsPerPage}`, {
            headers: { 'Authorization': `Bearer ${localStorage.getItem('lifeline_token')}` }
        });

        const banks = await response.json();
        
        // Render Banks
        renderBanks(banks); 

        const paginationHeader = response.headers.get("X-Pagination");
        if (paginationHeader) {
            const meta = JSON.parse(paginationHeader);
            renderPagination(meta);
        }

    } catch (error) {
        console.error("Failed to load banks", error);
    }
}

//  Render Pagination Buttons
function renderPagination(meta) {
    const controls = document.getElementById('paginationControls');
    controls.innerHTML = ''; 

    if (meta.totalPages <= 1) return; 
    const prevBtn = createPageBtn('<i class="fa-solid fa-chevron-right"></i>', meta.currentPage - 1, meta.currentPage > 1);
    controls.appendChild(prevBtn);

    for (let i = 1; i <= meta.totalPages; i++) {
        if (i === 1 || i === meta.totalPages || (i >= meta.currentPage - 1 && i <= meta.currentPage + 1)) {
            const btn = createPageBtn(i, i, true, i === meta.currentPage);
            controls.appendChild(btn);
        } else if (i === meta.currentPage - 2 || i === meta.currentPage + 2) {
            const span = document.createElement('span');
            span.className = "text-slate-400 px-2";
            span.innerText = "...";
            controls.appendChild(span);
        }
    }

    const nextBtn = createPageBtn('<i class="fa-solid fa-chevron-left"></i>', meta.currentPage + 1, meta.currentPage < meta.totalPages);
    controls.appendChild(nextBtn);
}

// Helper to create button
function createPageBtn(label, targetPage, isEnabled, isActive = false) {
    const btn = document.createElement('button');
    btn.innerHTML = label;
    
    let classes = "w-10 h-10 rounded-lg text-sm font-bold transition flex items-center justify-center ";
    if (isActive) {
        classes += "bg-morocco-red text-white shadow-lg";
    } else if (isEnabled) {
        classes += "bg-white text-slate-600 hover:bg-slate-100 border border-slate-200";
        btn.onclick = () => {
            currentPage = targetPage;
            loadBanks(targetPage);
            window.scrollTo({ top: 0, behavior: 'smooth' });
        };
    } else {
        classes += "bg-slate-50 text-slate-300 cursor-not-allowed";
    }
    
    btn.className = classes;
    return btn;
}

// Init Load
document.addEventListener("DOMContentLoaded", () => {
    loadBanks(currentPage);
});