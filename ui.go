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
				<button id="random-view-btn" class="px-3 py-2.5 sm:px-4 md:px-6 text-sm md:text-base rounded-md bg-emerald-600 text-white transition touch-manipulation">
					Random
				</button>
				<button id="list-view-btn" class="px-3 py-2.5 sm:px-4 md:px-6 text-sm md:text-base rounded-md text-slate-400 hover:text-slate-100 transition touch-manipulation">
					Declarations
				</button>
				<button id="labels-view-btn" class="px-3 py-2.5 sm:px-4 md:px-6 text-sm md:text-base rounded-md text-slate-400 hover:text-slate-100 transition touch-manipulation">
					Labels
				</button>
				<button id="env-view-btn" class="px-3 py-2.5 sm:px-4 md:px-6 text-sm md:text-base rounded-md text-slate-400 hover:text-slate-100 transition touch-manipulation hidden">
					Environment
				</button>
			</div>
		</div>

		<!-- Random View (Default) -->
		<div id="random-view" class="">
			<div class="bg-slate-900 rounded-lg shadow-lg p-6 mb-6 border border-slate-700">
				<h2 class="text-xl font-semibold mb-4">Random Declaration</h2>
				<button id="get-random-btn" class="px-6 py-3 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 transition touch-manipulation active:scale-95">
					Get Random Declaration
				</button>
			</div>
			<div id="random-declarations" class="space-y-4"></div>
		</div>

		<!-- List View -->
		<div id="list-view" class="hidden">
			<!-- Search -->
			<div class="bg-slate-900 rounded-lg shadow-lg p-4 mb-6 border border-slate-700">
				<input type="text" id="search-input" placeholder="Search declarations..."
					   class="w-full px-4 py-2 bg-slate-950 border border-slate-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-emerald-600 text-slate-100">
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
				<p class="text-slate-500">Try adjusting your search</p>
			</div>
		</div>
	</div>

		<!-- Environment View -->
		<div id="env-view" class="hidden">
			<div class="bg-slate-900 rounded-lg shadow-lg p-6 border border-slate-700">
				<h2 class="text-xl font-semibold mb-4">Environment Variables</h2>
				<div id="env-loading" class="text-center py-8">
					<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-emerald-600 mx-auto"></div>
					<p class="mt-4 text-slate-400">Loading environment variables...</p>
				</div>
				<div id="env-error" class="bg-red-950/60 border border-red-500/60 text-red-200 px-4 py-3 rounded-lg hidden">
					<span id="env-error-text"></span>
				</div>
				<div id="env-table-container" class="hidden overflow-x-auto">
					<table class="w-full">
						<thead class="bg-slate-800 border-b border-slate-700">
							<tr>
								<th class="px-6 py-3 text-left text-xs font-medium text-slate-300 uppercase tracking-wider">
									Variable Name
								</th>
								<th class="px-6 py-3 text-left text-xs font-medium text-slate-300 uppercase tracking-wider">
									Value
								</th>
							</tr>
						</thead>
						<tbody id="env-table" class="divide-y divide-slate-800">
						</tbody>
					</table>
				</div>
			</div>
		</div>

	<!-- Labels View -->
		<div id="labels-view" class="hidden">
			<div class="bg-slate-900 rounded-lg shadow-lg p-6 mb-6 border border-slate-700">
				<h2 class="text-xl font-semibold mb-4">Labels</h2>
				<div id="labels-loading" class="text-center py-8">
					<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-emerald-600 mx-auto"></div>
					<p class="mt-4 text-slate-400">Loading labels...</p>
				</div>
				<div id="labels-list" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-2 hidden"></div>
			</div>
			<div id="label-declarations" class="hidden">
				<div class="flex items-center justify-between mb-4">
					<h3 class="text-lg font-semibold text-slate-200">Declarations for <span id="selected-label-name" class="text-emerald-400"></span></h3>
					<button id="clear-label-selection" class="text-sm text-slate-400 hover:text-slate-200 transition">Show All Labels</button>
				</div>
				<div id="label-declarations-list" class="space-y-4"></div>
			</div>
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
let labels = [];
let currentView = 'random';
let selectedLabel = null;
let selectedLabelInLabelsView = null;
let envEnabled = false;

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
	loading: document.getElementById('loading'),
	errorMessage: document.getElementById('error-message'),
	errorText: document.getElementById('error-text'),
	resultsCount: document.getElementById('results-count'),
	declarationsTable: document.getElementById('declarations-table'),
	emptyState: document.getElementById('empty-state'),
	getRandomBtn: document.getElementById('get-random-btn'),
	randomDeclarations: document.getElementById('random-declarations'),
	labelsViewBtn: document.getElementById('labels-view-btn'),
	labelsView: document.getElementById('labels-view'),
	envViewBtn: document.getElementById('env-view-btn'),
	envView: document.getElementById('env-view'),
	envLoading: document.getElementById('env-loading'),
	envError: document.getElementById('env-error'),
	envErrorText: document.getElementById('env-error-text'),
	envTableContainer: document.getElementById('env-table-container'),
	envTable: document.getElementById('env-table'),
	labelsLoading: document.getElementById('labels-loading'),
	labelsList: document.getElementById('labels-list'),
	labelDeclarations: document.getElementById('label-declarations'),
	selectedLabelName: document.getElementById('selected-label-name'),
	clearLabelSelection: document.getElementById('clear-label-selection'),
	labelDeclarationsList: document.getElementById('label-declarations-list'),
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

async function loadLabels() {
	try {
		const data = await apiCall('/api/v1/labels');
		labels = data.labels || [];
		renderLabels();
	} catch (error) {
		console.error('Failed to load labels:', error);
	}
}

async function loadLabelDeclarations(label) {
	try {
		const data = await apiCall('/api/v1/declarations/label/' + encodeURIComponent(label));
		selectedLabelInLabelsView = label;
		renderLabelDeclarations(data.declarations || []);
	} catch (error) {
		console.error('Failed to load label declarations:', error);
	}
}

async function loadEnvironmentVariables() {
	try {
		elements.envLoading.classList.remove('hidden');
		elements.envError.classList.add('hidden');
		elements.envTableContainer.classList.add('hidden');
		
		const data = await apiCall('/api/v1/env');
		const envVars = data.environment || {};
		
		// Sort environment variables by name
		const sortedKeys = Object.keys(envVars).sort();
		
		elements.envTable.innerHTML = sortedKeys.map(key => {
			return '<tr class="hover:bg-slate-800 transition">' +
				'<td class="px-6 py-4 font-mono text-sm text-emerald-400">' + escapeHtml(key) + '</td>' +
				'<td class="px-6 py-4 font-mono text-sm text-slate-300 break-all">' + escapeHtml(envVars[key]) + '</td>' +
			'</tr>';
		}).join('');
		
		elements.envLoading.classList.add('hidden');
		elements.envTableContainer.classList.remove('hidden');
	} catch (error) {
		elements.envLoading.classList.add('hidden');
		elements.envErrorText.textContent = 'Failed to load environment variables: ' + error.message;
		elements.envError.classList.remove('hidden');
		console.error('Failed to load environment variables:', error);
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

		const response = await fetch('/api/v1/bible-esv/' + encodeURIComponent(reference));
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

	// Hide all views
	elements.listView.classList.add('hidden');
	elements.randomView.classList.add('hidden');
	elements.labelsView.classList.add('hidden');
	elements.envView.classList.add('hidden');

	// Reset all buttons
	const inactiveClass = 'px-3 py-2.5 sm:px-4 md:px-6 text-sm md:text-base rounded-md text-slate-400 hover:text-slate-100 transition touch-manipulation';
	const activeClass = 'px-3 py-2.5 sm:px-4 md:px-6 text-sm md:text-base rounded-md bg-emerald-600 text-white transition touch-manipulation';
	
	elements.listViewBtn.className = inactiveClass;
	elements.randomViewBtn.className = inactiveClass;
	elements.labelsViewBtn.className = inactiveClass;
	if (envEnabled) {
		elements.envViewBtn.className = inactiveClass;
	}

	if (view === 'list') {
		elements.listView.classList.remove('hidden');
		elements.listViewBtn.className = activeClass;
		loadDeclarations();
	} else if (view === 'labels') {
		elements.labelsView.classList.remove('hidden');
		elements.labelsViewBtn.className = activeClass;
		loadLabels();
	} else if (view === 'env' && envEnabled) {
		elements.envView.classList.remove('hidden');
		elements.envViewBtn.className = activeClass;
		loadEnvironmentVariables();
	} else {
		elements.randomView.classList.remove('hidden');
		elements.randomViewBtn.className = activeClass;
	}
}

function filterDeclarations() {
	const searchTerm = elements.searchInput.value.toLowerCase();
	let filtered = declarations;

	// Filter by label if one is selected
	if (selectedLabel) {
		filtered = filtered.filter(decl => {
			if (!decl.label) return false;
			const labels = decl.label.split(':').map(l => l.trim());
			return labels.includes(selectedLabel);
		});
	}

	// Filter by search term
	if (searchTerm) {
		filtered = filtered.filter(decl =>
			decl.text.toLowerCase().includes(searchTerm) ||
			decl.reference.toLowerCase().includes(searchTerm) ||
			(decl.label && decl.label.toLowerCase().includes(searchTerm))
		);
	}

	filteredDeclarations = filtered;
	renderDeclarations();
}

function renderDeclarations() {
	const searchTerm = elements.searchInput.value;

	// Show results count if filtering
	if (selectedLabel || searchTerm) {
		let message = 'Showing ' + filteredDeclarations.length + ' result' + (filteredDeclarations.length !== 1 ? 's' : '');
		if (selectedLabel) {
			message += ' for label "' + selectedLabel + '"';
		}
		if (searchTerm) {
			message += (selectedLabel ? ' and' : '') + ' search "' + searchTerm + '"';
		}
		elements.resultsCount.textContent = message;
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
		// Render labels as clickable badges
		let labelHtml = '';
		if (decl.label) {
			const labels = decl.label.split(':').map(l => l.trim()).filter(l => l);
			labelHtml = labels.map(label => {
				const isSelected = selectedLabel === label;
				const bgClass = isSelected ? 'bg-emerald-600 text-white' : 'bg-blue-500/20 text-blue-300 hover:bg-blue-500/30';
				return '<a href="#" onclick="event.preventDefault(); toggleLabelFilter(\'' +
					escapeHtml(label).replace(/'/g, "\\'") + '\');" ' +
					'class="inline-block ' + bgClass + ' text-xs px-2 py-1 rounded-full mr-2 cursor-pointer transition">' +
					escapeHtml(label) + '</a>';
			}).join('');
		}

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
		'</tr>';
	}).join('');
}

function renderLabels() {
	elements.labelsLoading.classList.add('hidden');
	elements.labelsList.classList.remove('hidden');
	elements.labelDeclarations.classList.add('hidden');
	
	elements.labelsList.innerHTML = labels.map(labelObj => {
		return '<a href="#" onclick="event.preventDefault(); loadLabelDeclarations(\'' +
			escapeHtml(labelObj.label).replace(/'/g, "\\'") + '\');" ' +
			'class="block py-2 px-3 bg-slate-800 hover:bg-slate-700 rounded transition text-sm">' +
			'<span class="text-slate-100">' + escapeHtml(labelObj.label) + '</span>' +
			'<span class="text-slate-400 text-xs ml-1.5">(' + labelObj.count + ')</span>' +
			'</a>';
	}).join('');
}

function renderLabelDeclarations(decls) {
	elements.labelsList.classList.add('hidden');
	elements.labelDeclarations.classList.remove('hidden');
	elements.selectedLabelName.textContent = selectedLabelInLabelsView;
	
	elements.labelDeclarationsList.innerHTML = decls.map(decl => {
		let labelHtml = '';
		if (decl.label) {
			const labels = decl.label.split(':').map(l => l.trim()).filter(l => l);
			labelHtml = labels.map(label => {
				return '<a href="#" onclick="event.preventDefault(); loadLabelDeclarations(\'' +
					escapeHtml(label).replace(/'/g, "\\'") + '\');" ' +
					'class="inline-block bg-emerald-500/20 text-emerald-300 hover:bg-emerald-500/30 text-xs px-2 py-1 rounded-full mr-2 cursor-pointer transition">' +
					escapeHtml(label) + '</a>';
			}).join('');
			labelHtml += '<br>';
		}
		return '<div class="bg-slate-900 rounded-lg shadow-lg p-6 mb-4 border border-slate-700">' +
			labelHtml +
			'<p class="text-slate-100 text-lg mb-2">' + escapeHtml(decl.text) + '</p>' +
			'<p class="text-slate-400">— <a href="#" onclick="event.preventDefault(); showBibleModal(\'' +
			escapeHtml(decl.reference).replace(/'/g, "\\'") + '\');" ' +
			'class="text-blue-400 hover:text-blue-300 hover:underline transition">' +
			escapeHtml(decl.reference) + '</a></p>' +
		'</div>';
	}).join('');
}

function displayRandomDeclaration(decl) {
	const div = document.createElement('div');
	div.className = 'bg-slate-900 rounded-lg shadow-lg p-6 border border-slate-700 animate-fade-in';

	// Render labels as clickable badges (like in the list view)
	let labelHtml = '';
	if (decl.label) {
		const labels = decl.label.split(':').map(l => l.trim()).filter(l => l);
		labelHtml = labels.map(label => {
			return '<a href="#" onclick="event.preventDefault(); filterByLabelAndSwitchToList(\'' +
				escapeHtml(label).replace(/'/g, "\\'") + '\');" ' +
				'class="inline-block bg-emerald-500/20 text-emerald-300 hover:bg-emerald-500/30 text-xs px-2 py-1 rounded-full mr-2 cursor-pointer transition">' +
				escapeHtml(label) + '</a>';
		}).join('');
		labelHtml += '<br>';
	}

	div.innerHTML = labelHtml +
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

function toggleLabelFilter(label) {
	// Toggle the label filter
	if (selectedLabel === label) {
		// Clicking the same label again clears the filter
		selectedLabel = null;
	} else {
		// Select the new label
		selectedLabel = label;
	}
	filterDeclarations();
}

function filterByLabelAndSwitchToList(label) {
	// Set the label filter and switch to list view
	selectedLabel = label;
	switchView('list');
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
	// Check for env query parameter
	const urlParams = new URLSearchParams(window.location.search);
	envEnabled = urlParams.get('env') === 'true';
	if (envEnabled) {
		elements.envViewBtn.classList.remove('hidden');
	}

	// View switching
	elements.listViewBtn.addEventListener('click', () => switchView('list'));
	elements.randomViewBtn.addEventListener('click', () => switchView('random'));
	elements.labelsViewBtn.addEventListener('click', () => switchView('labels'));
	if (envEnabled) {
		elements.envViewBtn.addEventListener('click', () => switchView('env'));
	}

	// Search
	elements.searchInput.addEventListener('input', filterDeclarations);

	// Random declaration
	elements.getRandomBtn.addEventListener('click', getRandomDeclaration);

	// Labels view
	elements.clearLabelSelection.addEventListener('click', () => {
		selectedLabelInLabelsView = null;
		renderLabels();
	});

	// Bible Modal
	elements.bibleModalClose.addEventListener('click', hideBibleModal);
	elements.bibleModal.addEventListener('click', (e) => {
		if (e.target === elements.bibleModal) hideBibleModal();
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

/* Mobile touch optimization */
@media (hover: none) and (pointer: coarse) {
	/* Increase tap target size for mobile */
	a, button {
		min-height: 44px;
		display: inline-flex;
		align-items: center;
		justify-content: center;
	}
	
	/* Remove hover effects on touch devices */
	button:hover {
		transform: none;
	}
	
	/* Add active state feedback for better touch UX */
	button:active {
		opacity: 0.8;
	}
}

/* Ensure buttons don't wrap on mobile even with 4 buttons */
@media (max-width: 640px) {
	#random-view-btn, #list-view-btn, #labels-view-btn, #env-view-btn {
		font-size: 0.8125rem; /* 13px */
		padding-left: 0.625rem; /* 10px */
		padding-right: 0.625rem; /* 10px */
	}
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
