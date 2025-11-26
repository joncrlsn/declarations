package main

import (
	"html/template"
	"net/http"
)

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Declarations in Christ</title>
	<script src="https://cdn.tailwindcss.com"></script>
</head>
<body id="app-body" class="min-h-screen bg-slate-950 text-slate-50">
	<div class="container mx-auto px-4 py-8">
		<!-- Header -->
		<div class="flex flex-col md:flex-row md:items-center md:justify-between mb-8 gap-4">
			<div class="text-center md:text-left">
				<h1 class="text-3xl font-bold mb-2">Declarations in Christ</h1>
				<p class="text-slate-300">Biblical declarations of who you are in Christ</p>
				<div class="mt-4 flex items-center justify-center md:justify-start text-sm text-slate-400">
					<span id="status-indicator" class="inline-block w-2 h-2 bg-slate-500 rounded-full mr-2"></span>
					<span id="status-text">Checking API...</span>
					<span class="mx-2">•</span>
					<span id="total-count">0</span> declarations
				</div>
			</div>
		</div>

		<!-- View Toggle -->
		<div class="flex justify-center mb-6">
			<div class="bg-slate-900 rounded-lg p-1 shadow-sm border border-slate-700">
				<button id="list-view-btn" class="px-4 py-2 rounded-md text-slate-200 hover:text-white">List View</button>
				<button id="random-view-btn" class="px-4 py-2 rounded-md bg-emerald-500 text-white">Random View</button>
			</div>
		</div>

		<!-- List View -->
		<div id="list-view" class="hidden">
			<!-- Search and Add -->
			<div class="bg-slate-900 rounded-lg shadow-sm p-4 mb-6 border border-slate-700">
				<div class="flex flex-col md:flex-row gap-4">
					<input type="text" id="search-input" placeholder="Search declarations..."
						   class="flex-1 px-3 py-2 border border-slate-600 bg-slate-950 rounded-md focus:outline-none focus:ring-2 focus:ring-emerald-500 text-slate-100 placeholder:text-slate-500">
					<button id="add-btn" class="px-4 py-2 bg-emerald-500 text-white rounded-md hover:bg-emerald-600">
						Add Declaration
					</button>
				</div>
			</div>

			<!-- Loading State -->
			<div id="loading" class="text-center py-12">
				<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-emerald-500 mx-auto"></div>
				<p class="mt-2 text-slate-300">Loading declarations...</p>
			</div>

			<!-- Error Message -->
			<div id="error-message" class="bg-red-950/60 border border-red-500/60 text-red-200 px-4 py-3 rounded-md mb-6 hidden">
                <span id="error-text"></span>
            </div>

			<!-- Results Count -->
			<div id="results-count" class="text-sm text-slate-300 mb-4 hidden"></div>

			<!-- Declarations Table -->
			<div class="bg-slate-900 rounded-lg shadow-sm overflow-hidden border border-slate-700">
				<table class="w-full">
					<thead class="bg-slate-800">
						<tr>
							<th class="px-4 py-3 text-left text-xs font-medium text-slate-300 uppercase">Declaration</th>
							<th class="px-4 py-3 text-left text-xs font-medium text-slate-300 uppercase">Reference</th>
							<th class="px-4 py-3 text-right text-xs font-medium text-slate-300 uppercase">Actions</th>
                        </tr>
                    </thead>
                    <tbody id="declarations-table" class="divide-y divide-gray-200">
                    </tbody>
                </table>
            </div>

			<!-- Empty State -->
			<div id="empty-state" class="text-center py-12 hidden">
				<div class="text-slate-500 mb-4">
					<svg class="mx-auto h-12 w-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                    </svg>
                </div>
				<h3 class="text-lg font-medium text-slate-100 mb-2">No declarations found</h3>
				<p class="text-slate-300">Try adjusting your search or add a new declaration</p>
			</div>
		</div>

		<!-- Random View -->
		<div id="random-view" class="">
			<div class="bg-slate-900 rounded-lg shadow-sm p-6 mb-6 border border-slate-700">
				<h2 class="text-xl font-semibold text-slate-100 mb-2">Random Declarations</h2>
				<p class="text-slate-300 mb-4">Click the button to see more.</p>
				<button id="get-random-btn" class="px-6 py-3 bg-emerald-500 text-white rounded-md hover:bg-emerald-600">
					Get Random Declaration
				</button>
			</div>
			<div id="random-declarations" class="space-y-4"></div>
		</div>
	</div>

	<!-- Add/Edit Modal -->
	<div id="modal" class="fixed inset-0 bg-black bg-opacity-60 hidden items-center justify-center">
		<div class="bg-slate-900 rounded-lg p-6 w-full max-w-md mx-4 border border-slate-700">
			<h3 id="modal-title" class="text-lg font-semibold text-slate-100 mb-4">Add Declaration</h3>
			<form id="declaration-form">
				<div class="mb-4">
					<label class="block text-sm font-medium text-slate-200 mb-2">Label (optional)</label>
					<input type="text" id="label-input" placeholder="e.g., Promise, Blessing"
						   class="w-full px-3 py-2 border border-slate-600 bg-slate-950 rounded-md focus:outline-none focus:ring-2 focus:ring-emerald-500 text-slate-100 placeholder:text-slate-500">
				</div>
				<div class="mb-4">
					<label class="block text-sm font-medium text-slate-200 mb-2">Declaration Text</label>
					<textarea id="text-input" rows="3" placeholder="Enter the declaration text..."
							  class="w-full px-3 py-2 border border-slate-600 bg-slate-950 rounded-md focus:outline-none focus:ring-2 focus:ring-emerald-500 text-slate-100 placeholder:text-slate-500" required></textarea>
				</div>
				<div class="mb-6">
					<label class="block text-sm font-medium text-slate-200 mb-2">Bible Reference</label>
					<input type="text" id="reference-input" placeholder="e.g., John 3:16"
						   class="w-full px-3 py-2 border border-slate-600 bg-slate-950 rounded-md focus:outline-none focus:ring-2 focus:ring-emerald-500 text-slate-100 placeholder:text-slate-500" required>
				</div>
				<div class="flex gap-3">
					<button type="submit" class="flex-1 px-4 py-2 bg-emerald-500 text-white rounded-md hover:bg-emerald-600">
						Save
					</button>
					<button type="button" id="cancel-btn" class="flex-1 px-4 py-2 border border-slate-600 text-slate-200 rounded-md hover:bg-slate-800">
						Cancel
					</button>
				</div>
			</form>
		</div>
	</div>

<script>
// Global state
let declarations = [];
let filteredDeclarations = [];
let currentView = 'random';
let editingId = null;

// DOM elements
const elements = {
    // Status
    statusIndicator: document.getElementById('status-indicator'),
    statusText: document.getElementById('status-text'),
    totalCount: document.getElementById('total-count'),

    // Views
    listView: document.getElementById('list-view'),
    randomView: document.getElementById('random-view'),
    listViewBtn: document.getElementById('list-view-btn'),
    randomViewBtn: document.getElementById('random-view-btn'),

    // List view
    searchInput: document.getElementById('search-input'),
    addBtn: document.getElementById('add-btn'),
    loading: document.getElementById('loading'),
    errorMessage: document.getElementById('error-message'),
    errorText: document.getElementById('error-text'),
    resultsCount: document.getElementById('results-count'),
    declarationsTable: document.getElementById('declarations-table'),
    emptyState: document.getElementById('empty-state'),

    // Random view
    getRandomBtn: document.getElementById('get-random-btn'),
    randomDeclarations: document.getElementById('random-declarations'),

    // Modal
    modal: document.getElementById('modal'),
    modalTitle: document.getElementById('modal-title'),
    declarationForm: document.getElementById('declaration-form'),
    labelInput: document.getElementById('label-input'),
    textInput: document.getElementById('text-input'),
    referenceInput: document.getElementById('reference-input'),
    cancelBtn: document.getElementById('cancel-btn')
};

// API functions
async function apiCall(endpoint, options = {}) {
    const response = await fetch(endpoint, {
        headers: { 'Content-Type': 'application/json', ...options.headers },
        ...options
    });
    if (!response.ok) {
        throw new Error('API request failed: ' + response.status);
    }
    return response.json();
}

async function loadDeclarations() {
    try {
        showLoading(true);
        const data = await apiCall('/api/v1/declarations');
        declarations = data.declarations || [];
        filterDeclarations();
        showError(null);
    } catch (error) {
        showError('Failed to load declarations: ' + error.message);
        console.error('Load error:', error);
    } finally {
        showLoading(false);
    }
}

async function checkApiHealth() {
    try {
        const health = await apiCall('/api/v1/health');
        elements.statusIndicator.className = 'inline-block w-2 h-2 bg-green-500 rounded-full mr-2';
        elements.statusText.textContent = 'API Connected';
        elements.totalCount.textContent = health.declarations_count;
    } catch (error) {
        elements.statusIndicator.className = 'inline-block w-2 h-2 bg-red-500 rounded-full mr-2';
        elements.statusText.textContent = 'API Error';
        elements.totalCount.textContent = '0';
    }
}

async function saveDeclaration(declarationData) {
    const endpoint = editingId ? '/api/v1/declarations/' + editingId : '/api/v1/declarations';
    const method = editingId ? 'PUT' : 'POST';

    await apiCall(endpoint, {
        method: method,
        body: JSON.stringify(declarationData)
    });
}

async function deleteDeclaration(id) {
    if (!confirm('Are you sure you want to delete this declaration?')) return;

    try {
        await apiCall('/api/v1/declarations/' + id, { method: 'DELETE' });
        await loadDeclarations();
        await checkApiHealth();
    } catch (error) {
        showError('Failed to delete declaration: ' + error.message);
    }
}

async function getRandomDeclaration() {
    try {
        const declaration = await apiCall('/api/v1/declarations/random');
        displayRandomDeclaration(declaration);
    } catch (error) {
        showError('Failed to get random declaration: ' + error.message);
    }
}

// UI functions
function showLoading(show) {
    elements.loading.classList.toggle('hidden', !show);
}

function showError(message) {
    if (message) {
        elements.errorText.textContent = message;
        elements.errorMessage.classList.remove('hidden');
    } else {
        elements.errorMessage.classList.add('hidden');
    }
}

function switchView(view) {
    currentView = view;

    if (view === 'list') {
        elements.listView.classList.remove('hidden');
        elements.randomView.classList.add('hidden');
        elements.listViewBtn.className = 'px-4 py-2 rounded-md bg-emerald-500 text-white';
        elements.randomViewBtn.className = 'px-4 py-2 rounded-md text-slate-200 hover:text-white';
    } else {
        elements.listView.classList.add('hidden');
        elements.randomView.classList.remove('hidden');
        elements.listViewBtn.className = 'px-4 py-2 rounded-md text-slate-200 hover:text-white';
        elements.randomViewBtn.className = 'px-4 py-2 rounded-md bg-emerald-500 text-white';
    }
}

function filterDeclarations() {
    const searchTerm = elements.searchInput.value.toLowerCase();
    filteredDeclarations = declarations.filter(decl =>
        decl.text.toLowerCase().includes(searchTerm) ||
        decl.reference.toLowerCase().includes(searchTerm) ||
        (decl.label && decl.label.toLowerCase().includes(searchTerm))
    );
    renderDeclarations();
}

function renderDeclarations() {
    const searchTerm = elements.searchInput.value;

    // Update results count
    if (searchTerm) {
        elements.resultsCount.textContent = 'Showing ' + filteredDeclarations.length + ' results for "' + searchTerm + '"';
        elements.resultsCount.classList.remove('hidden');
    } else {
        elements.resultsCount.classList.add('hidden');
    }

    // Show empty state if no results
    if (filteredDeclarations.length === 0) {
        elements.declarationsTable.innerHTML = '';
        elements.emptyState.classList.remove('hidden');
        return;
    }

    elements.emptyState.classList.add('hidden');

    // Render table rows - compact format, no ID shown
    elements.declarationsTable.innerHTML = filteredDeclarations.map(decl => {
        const labelHtml = decl.label ?
            '<span class="inline-block bg-blue-100 text-blue-800 text-xs px-2 py-1 rounded-full mr-2">' +
            highlightText(decl.label, searchTerm) + '</span>' : '';

        // In dark mode, default row is dark background with light text.
        // When the row is highlighted (hovered), switch to light background and dark text.
        return '<tr class="group bg-slate-950 text-slate-100 hover:bg-slate-100 hover:text-slate-900">' +
            '<td class="px-4 py-3 group-hover:text-slate-900">' + labelHtml + highlightText(decl.text, searchTerm) + '</td>' +
            '<td class="px-4 py-3 text-sm text-slate-300 group-hover:text-slate-900">' + highlightText(decl.reference, searchTerm) + '</td>' +
            '<td class="px-4 py-3 text-right group-hover:text-slate-900">' +
                '<button onclick="editDeclaration(' + decl.id + ')" class="text-blue-400 hover:text-blue-600 mr-2">Edit</button>' +
                '<button onclick="deleteDeclaration(' + decl.id + ')" class="text-red-400 hover:text-red-600">Delete</button>' +
            '</td>' +
        '</tr>';
    }).join('');
}

function highlightText(text, searchTerm) {
    if (!searchTerm) return text;
    const regex = new RegExp('(' + searchTerm.replace(/[.*+?^${}()|[\]\\]/g, '\\$&') + ')', 'gi');
    return text.replace(regex, '<mark>$1</mark>');
}

function displayRandomDeclaration(decl) {
    const div = document.createElement('div');
    div.className = 'bg-slate-900 rounded-lg shadow-sm p-6 border border-slate-700';

    const labelHtml = decl.label ?
        '<span class="inline-block bg-emerald-500/20 text-emerald-300 text-xs px-2 py-1 rounded-full mb-2">' + decl.label + '</span><br>' : '';

    div.innerHTML = labelHtml +
        '<p class="text-slate-50 text-lg mb-2">' + decl.text + '</p>' +
        '<p class="text-slate-300">— ' + decl.reference + '</p>';

    // Add new random declarations above previous ones
    if (elements.randomDeclarations.firstChild) {
        elements.randomDeclarations.insertBefore(div, elements.randomDeclarations.firstChild);
    } else {
        elements.randomDeclarations.appendChild(div);
    }
}

function showModal(title, declaration = null) {
    elements.modalTitle.textContent = title;
    elements.modal.classList.remove('hidden');
    elements.modal.classList.add('flex');

    if (declaration) {
        editingId = declaration.id;
        elements.labelInput.value = declaration.label || '';
        elements.textInput.value = declaration.text;
        elements.referenceInput.value = declaration.reference;
    } else {
        editingId = null;
        elements.declarationForm.reset();
    }

    elements.textInput.focus();
}

function hideModal() {
    elements.modal.classList.add('hidden');
    elements.modal.classList.remove('flex');
    editingId = null;
}

function editDeclaration(id) {
    const declaration = declarations.find(d => d.id === id);
    if (declaration) {
        showModal('Edit Declaration', declaration);
    }
}

// Event listeners
document.addEventListener('DOMContentLoaded', () => {
    // View switching
    elements.listViewBtn.addEventListener('click', () => switchView('list'));
    elements.randomViewBtn.addEventListener('click', () => switchView('random'));

    // Search
    elements.searchInput.addEventListener('input', filterDeclarations);

    // Add button
    elements.addBtn.addEventListener('click', () => showModal('Add Declaration'));

    // Random declaration
    elements.getRandomBtn.addEventListener('click', getRandomDeclaration);

    // Modal
    elements.cancelBtn.addEventListener('click', hideModal);
    elements.modal.addEventListener('click', (e) => {
        if (e.target === elements.modal) hideModal();
    });

    // Form submission
    elements.declarationForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const declarationData = {
            label: elements.labelInput.value.trim() || null,
            text: elements.textInput.value.trim(),
            reference: elements.referenceInput.value.trim()
        };

        try {
            await saveDeclaration(declarationData);
            hideModal();
            await loadDeclarations();
            await checkApiHealth();
        } catch (error) {
            showError('Failed to save declaration: ' + error.message);
        }
    });

    // Initialize
    loadDeclarations();
    checkApiHealth();
    switchView('random');
    getRandomDeclaration();
});
</script>
</body>
</html>
`

// ServeUI handles serving the web interface
func ServeUI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.New("ui").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}
