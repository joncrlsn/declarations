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
<body class="min-h-screen bg-slate-950 text-slate-50">
	<div class="container mx-auto px-4 py-8 max-w-6xl">
		<!-- Header -->
		<div class="mb-8">
			<h1 class="text-4xl font-bold text-center mb-2">Declarations in Christ</h1>
			<p class="text-slate-400 text-center">Biblical declarations of who you are in Christ</p>
			<div class="mt-4 flex items-center justify-center text-sm text-slate-400">
				<span id="status-indicator" class="inline-block w-2 h-2 bg-slate-500 rounded-full mr-2"></span>
				<span id="status-text">Checking API...</span>
				<span class="mx-2">•</span>
				<span id="total-count">0</span>&nbsp;declarations
			</div>
		</div>

		<!-- View Toggle -->
		<div class="flex justify-center mb-6">
			<div class="bg-slate-900 rounded-lg p-1 border border-slate-700">
				<button id="list-view-btn" class="px-6 py-2 rounded-md text-slate-400 hover:text-slate-100 transition">
					List View
				</button>
				<button id="random-view-btn" class="px-6 py-2 rounded-md bg-emerald-600 text-white transition">
					Random View
				</button>
			</div>
		</div>

		<!-- Random View (Default) -->
		<div id="random-view" class="">
			<div class="bg-slate-900 rounded-lg shadow-lg p-6 mb-6 border border-slate-700">
				<h2 class="text-xl font-semibold mb-4">Random Declaration</h2>
				<button id="get-random-btn" class="px-6 py-3 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 transition">
					Get Random Declaration
				</button>
			</div>
			<div id="random-declarations" class="space-y-4"></div>
		</div>

		<!-- List View -->
		<div id="list-view" class="hidden">
			<!-- Search and Add -->
			<div class="bg-slate-900 rounded-lg shadow-lg p-4 mb-6 border border-slate-700">
				<div class="flex flex-col md:flex-row gap-4">
					<input type="text" id="search-input" placeholder="Search declarations..."
						   class="flex-1 px-4 py-2 bg-slate-950 border border-slate-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-emerald-600 text-slate-100">
					<button id="add-btn" class="px-6 py-2 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 transition">
						Add Declaration
					</button>
				</div>
			</div>

			<!-- Loading State -->
			<div id="loading" class="text-center py-12">
				<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-emerald-600 mx-auto"></div>
				<p class="mt-4 text-slate-400">Loading declarations...</p>
			</div>

			<!-- Error Message -->
			<div id="error-message" class="bg-red-950/60 border border-red-500/60 text-red-200 px-4 py-3 rounded-lg mb-6 hidden">
				<span id="error-text"></span>
			</div>

			<!-- Results Count -->
			<div id="results-count" class="text-sm text-slate-400 mb-4 hidden"></div>

			<!-- Declarations Table -->
			<div class="bg-slate-900 rounded-lg shadow-lg overflow-hidden border border-slate-700">
				<table class="w-full">
					<thead class="bg-slate-800 border-b border-slate-700">
						<tr>
							<th class="px-6 py-3 text-left text-xs font-medium text-slate-300 uppercase tracking-wider">
								Declaration
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-slate-300 uppercase tracking-wider">
								Reference
							</th>
							<th class="px-6 py-3 text-right text-xs font-medium text-slate-300 uppercase tracking-wider">
								Actions
							</th>
						</tr>
					</thead>
					<tbody id="declarations-table" class="divide-y divide-slate-800">
					</tbody>
				</table>
			</div>

			<!-- Empty State -->
			<div id="empty-state" class="text-center py-12 hidden">
				<div class="text-slate-600 mb-4">
					<svg class="mx-auto h-16 w-16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
							  d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z">
						</path>
					</svg>
				</div>
				<h3 class="text-lg font-medium text-slate-300 mb-2">No declarations found</h3>
				<p class="text-slate-500">Try adjusting your search or add a new declaration</p>
			</div>
		</div>
	</div>

	<!-- Add/Edit Modal -->
	<div id="modal" class="fixed inset-0 bg-black bg-opacity-70 hidden items-center justify-center z-50">
		<div class="bg-slate-900 rounded-lg p-6 w-full max-w-md mx-4 border border-slate-700 shadow-2xl">
			<h3 id="modal-title" class="text-xl font-semibold mb-6">Add Declaration</h3>
			<form id="declaration-form">
				<div class="mb-4">
					<label class="block text-sm font-medium text-slate-300 mb-2">Label (optional)</label>
					<input type="text" id="label-input" placeholder="e.g., Promise, Blessing"
						   class="w-full px-4 py-2 bg-slate-950 border border-slate-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-emerald-600 text-slate-100">
				</div>
				<div class="mb-4">
					<label class="block text-sm font-medium text-slate-300 mb-2">Declaration Text</label>
					<textarea id="text-input" rows="3" placeholder="Enter the declaration text..." required
							  class="w-full px-4 py-2 bg-slate-950 border border-slate-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-emerald-600 text-slate-100"></textarea>
				</div>
				<div class="mb-6">
					<label class="block text-sm font-medium text-slate-300 mb-2">Bible Reference</label>
					<input type="text" id="reference-input" placeholder="e.g., John 3:16" required
						   class="w-full px-4 py-2 bg-slate-950 border border-slate-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-emerald-600 text-slate-100">
				</div>
				<div class="flex gap-3">
					<button type="submit" class="flex-1 px-4 py-2 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 transition">
						Save
					</button>
					<button type="button" id="cancel-btn" class="flex-1 px-4 py-2 bg-slate-800 text-slate-300 rounded-lg hover:bg-slate-700 transition">
						Cancel
					</button>
				</div>
			</form>
		</div>
	</div>

	<!-- Bible Text Modal -->
	<div id="bible-modal" class="fixed inset-0 bg-black bg-opacity-70 hidden items-center justify-center z-50">
		<div class="bg-slate-900 rounded-lg p-6 w-full max-w-2xl mx-4 border border-slate-700 shadow-2xl max-h-[80vh] overflow-y-auto">
			<div class="flex justify-end mb-4">
				<button id="bible-modal-close" class="text-slate-500 hover:text-slate-300 transition">
					<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
					</svg>
				</button>
			</div>
			<div id="bible-loading" class="text-center py-8 hidden">
				<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-emerald-600 mx-auto"></div>
				<p class="mt-4 text-slate-400">Loading Bible text...</p>
			</div>
			<div id="bible-text-content" class="text-slate-200 whitespace-pre-wrap leading-relaxed hidden"></div>
			<div id="bible-error" class="bg-red-950/60 border border-red-500/60 text-red-200 px-4 py-3 rounded-lg hidden"></div>
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
	statusIndicator: document.getElementById('status-indicator'),
	statusText: document.getElementById('status-text'),
	totalCount: document.getElementById('total-count'),
	listView: document.getElementById('list-view'),
	randomView: document.getElementById('random-view'),
	listViewBtn: document.getElementById('list-view-btn'),
	randomViewBtn: document.getElementById('random-view-btn'),
	searchInput: document.getElementById('search-input'),
	addBtn: document.getElementById('add-btn'),
	loading: document.getElementById('loading'),
	errorMessage: document.getElementById('error-message'),
	errorText: document.getElementById('error-text'),
	resultsCount: document.getElementById('results-count'),
	declarationsTable: document.getElementById('declarations-table'),
	emptyState: document.getElementById('empty-state'),
	getRandomBtn: document.getElementById('get-random-btn'),
	randomDeclarations: document.getElementById('random-declarations'),
	modal: document.getElementById('modal'),
	modalTitle: document.getElementById('modal-title'),
	declarationForm: document.getElementById('declaration-form'),
	labelInput: document.getElementById('label-input'),
	textInput: document.getElementById('text-input'),
	referenceInput: document.getElementById('reference-input'),
	cancelBtn: document.getElementById('cancel-btn'),
	bibleModal: document.getElementById('bible-modal'),
	bibleModalClose: document.getElementById('bible-modal-close'),
	bibleTextContent: document.getElementById('bible-text-content'),
	bibleLoading: document.getElementById('bible-loading'),
	bibleError: document.getElementById('bible-error')
};

// API functions
async function apiCall(endpoint, options = {}) {
	const response = await fetch(endpoint, {
		headers: { 'Content-Type': 'application/json', ...options.headers },
		...options
	});
	if (!response.ok && response.status !== 204) {
		throw new Error('API request failed: ' + response.status);
	}
	if (response.status === 204) {
		return null;
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

async function fetchBibleText(reference) {
	try {
		elements.bibleLoading.classList.remove('hidden');
		elements.bibleTextContent.classList.add('hidden');
		elements.bibleError.classList.add('hidden');

		const response = await fetch('/api/v1/bible/text?q=' + encodeURIComponent(reference));
		if (!response.ok) {
			throw new Error('Failed to fetch Bible text');
		}

		const data = await response.json();

		let text = '';
		if (data.passages && data.passages.length > 0) {
			text = data.passages.join('\n\n');
		} else {
			text = 'No text found for this reference.';
		}

		elements.bibleTextContent.textContent = text;
		elements.bibleTextContent.classList.remove('hidden');
		elements.bibleLoading.classList.add('hidden');
	} catch (error) {
		elements.bibleLoading.classList.add('hidden');
		elements.bibleError.textContent = 'Failed to load Bible text: ' + error.message;
		elements.bibleError.classList.remove('hidden');
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
		elements.listViewBtn.className = 'px-6 py-2 rounded-md bg-emerald-600 text-white transition';
		elements.randomViewBtn.className = 'px-6 py-2 rounded-md text-slate-400 hover:text-slate-100 transition';
		loadDeclarations();
	} else {
		elements.listView.classList.add('hidden');
		elements.randomView.classList.remove('hidden');
		elements.listViewBtn.className = 'px-6 py-2 rounded-md text-slate-400 hover:text-slate-100 transition';
		elements.randomViewBtn.className = 'px-6 py-2 rounded-md bg-emerald-600 text-white transition';
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

	if (searchTerm) {
		elements.resultsCount.textContent = 'Showing ' + filteredDeclarations.length + ' results for "' + searchTerm + '"';
		elements.resultsCount.classList.remove('hidden');
	} else {
		elements.resultsCount.classList.add('hidden');
	}

	if (filteredDeclarations.length === 0) {
		elements.declarationsTable.innerHTML = '';
		elements.emptyState.classList.remove('hidden');
		return;
	}

	elements.emptyState.classList.add('hidden');

	elements.declarationsTable.innerHTML = filteredDeclarations.map(decl => {
		const labelHtml = decl.label ?
			'<span class="inline-block bg-blue-500/20 text-blue-300 text-xs px-2 py-1 rounded-full mr-2">' +
			escapeHtml(decl.label) + '</span>' : '';

		return '<tr class="hover:bg-slate-800 transition">' +
			'<td class="px-6 py-4">' +
				labelHtml + escapeHtml(decl.text) +
			'</td>' +
			'<td class="px-6 py-4 text-sm">' +
				'<a href="#" onclick="event.preventDefault(); showBibleModal(\'' +
				escapeHtml(decl.reference).replace(/'/g, "\\'") + '\');" ' +
				'class="text-blue-400 hover:text-blue-300 hover:underline transition">' +
				escapeHtml(decl.reference) + '</a>' +
			'</td>' +
			'<td class="px-6 py-4 text-right space-x-2">' +
				'<button onclick="editDeclaration(' + decl.id + ')" ' +
				'class="text-blue-400 hover:text-blue-300 transition">Edit</button>' +
				'<button onclick="deleteDeclaration(' + decl.id + ')" ' +
				'class="text-red-400 hover:text-red-300 transition">Delete</button>' +
			'</td>' +
		'</tr>';
	}).join('');
}

function displayRandomDeclaration(decl) {
	const div = document.createElement('div');
	div.className = 'bg-slate-900 rounded-lg shadow-lg p-6 border border-slate-700 animate-fade-in';

	const labelHtml = decl.label ?
		'<span class="inline-block bg-emerald-500/20 text-emerald-300 text-xs px-2 py-1 rounded-full mb-2">' +
		escapeHtml(decl.label) + '</span><br>' : '';

	//div.innerHTML = labelHtml +
	div.innerHTML =
		'<p class="text-slate-100 text-lg mb-3 leading-relaxed">' + escapeHtml(decl.text) + '</p>' +
		'<p class="text-slate-400">— <a href="#" onclick="event.preventDefault(); showBibleModal(\'' +
		escapeHtml(decl.reference).replace(/'/g, "\\'") + '\');" ' +
		'class="text-blue-400 hover:text-blue-300 hover:underline transition">' +
		escapeHtml(decl.reference) + '</a></p>';

	// Add new declarations above previous ones
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

function showBibleModal(reference) {
	elements.bibleModal.classList.remove('hidden');
	elements.bibleModal.classList.add('flex');
	fetchBibleText(reference);
}

function hideBibleModal() {
	elements.bibleModal.classList.add('hidden');
	elements.bibleModal.classList.remove('flex');
	elements.bibleTextContent.textContent = '';
	elements.bibleError.classList.add('hidden');
}

function escapeHtml(text) {
	const div = document.createElement('div');
	div.textContent = text;
	return div.innerHTML;
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

	// Bible Modal
	elements.bibleModalClose.addEventListener('click', hideBibleModal);
	elements.bibleModal.addEventListener('click', (e) => {
		if (e.target === elements.bibleModal) hideBibleModal();
	});

	// Form submission
	elements.declarationForm.addEventListener('submit', async (e) => {
		e.preventDefault();

		const declarationData = {
			label: elements.labelInput.value.trim() || '',
			text: elements.textInput.value.trim(),
			reference: elements.referenceInput.value.trim()
		};

		try {
			await saveDeclaration(declarationData);
			hideModal();
			if (currentView === 'list') {
				await loadDeclarations();
			}
			await checkApiHealth();
		} catch (error) {
			showError('Failed to save declaration: ' + error.message);
		}
	});

	// Initialize
	checkApiHealth();
	switchView('random');
	getRandomDeclaration(); // Pre-load first random declaration
});
</script>

<style>
@keyframes fade-in {
	from { opacity: 0; transform: translateY(-10px); }
	to { opacity: 1; transform: translateY(0); }
}
.animate-fade-in {
	animation: fade-in 0.3s ease-out;
}
</style>

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

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}
