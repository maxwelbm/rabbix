// Global variables
let currentExecution = null;
let eventSource = null;
let statusChart = null;
let timingChart = null;
let draggedElement = null;
let executionOrder = [];

// Initialize the application
document.addEventListener('DOMContentLoaded', function() {
    updateActiveExecutions();
    initializeCharts();
    initializeDragAndDrop();
    updateExecutionOrder();
});

// Tab management
function showTab(tabName) {
    // Remove active class from all tabs and panels
    document.querySelectorAll('.tab').forEach(tab => tab.classList.remove('active'));
    document.querySelectorAll('.tab-panel').forEach(panel => panel.classList.remove('active'));
    
    // Add active class to selected tab and panel
    document.querySelector('.tab[onclick="showTab(\'' + tabName + '\')"]').classList.add('active');
    document.getElementById(tabName + '-tab').classList.add('active');
}

// Test selection functions
function selectAllTests() {
    document.querySelectorAll('.test-checkbox').forEach(checkbox => {
        checkbox.checked = true;
    });
}

function clearSelection() {
    document.querySelectorAll('.test-checkbox').forEach(checkbox => {
        checkbox.checked = false;
    });
}

function getSelectedTests() {
    // Return tests in execution order
    if (executionOrder.length > 0) {
        return executionOrder;
    }
    return Array.from(document.querySelectorAll('.test-checkbox:checked'))
        .map(checkbox => checkbox.value);
}

// Run single test
async function runSingleTest(testName) {
    const btn = event.target;
    const originalText = btn.textContent;
    btn.textContent = '⏳';
    btn.disabled = true;

    try {
        const response = await fetch('/api/run/' + testName, {
            method: 'POST'
        });
        
        const result = await response.json();
        
        // Add to logs
        addLogEntry('info', 'Teste individual executado: ' + testName);
        addLogEntry(result.status, testName + ': ' + result.status.toUpperCase() + ' (' + result.duration_ms + 'ms)');
        
        // Add to results
        addSingleResult(result);
        
    } catch (error) {
        addLogEntry('error', 'Erro ao executar ' + testName + ': ' + error.message);
    } finally {
        btn.textContent = originalText;
        btn.disabled = false;
    }
}

// Execute batch
async function executeBatch() {
    const selectedTests = getSelectedTests();
    
    if (selectedTests.length === 0) {
        alert('Selecione pelo menos um teste para executar em lote.');
        return;
    }
    
    const concurrency = parseInt(document.getElementById('concurrency').value);
    const delay = parseInt(document.getElementById('delay').value);
    
    clearLogs();
    clearResults();
    showTab('logs');
    
    // Log execution order
    addLogEntry('info', '🚀 Iniciando execução em lote');
    addLogEntry('info', '📋 Ordem de execução definida:');
    selectedTests.forEach((test, index) => {
        addLogEntry('info', `  ${index + 1}. ${test}`);
    });
    addLogEntry('info', `⚙️ Configurações: Concorrência=${concurrency}, Delay=${delay}ms`);
    
    try {
        const response = await fetch('/api/batch', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                tests: selectedTests,
                concurrency: concurrency,
                delay: delay
            })
        });
        
        const result = await response.json();
        currentExecution = result.execution_id;
        
        // Start listening to logs
        startLogStream(currentExecution);
        
        // Start polling for execution status
        pollExecutionStatus(currentExecution);
        
        updateActiveExecutions();
        
    } catch (error) {
        addLogEntry('error', 'Erro ao iniciar execução em lote: ' + error.message);
    }
}

// Start log stream
function startLogStream(executionId) {
    if (eventSource) {
        eventSource.close();
    }
    
    eventSource = new EventSource('/api/logs/' + executionId);
    
    eventSource.onmessage = function(event) {
        const logMessage = JSON.parse(event.data);
        addLogEntry(logMessage.level, logMessage.message);
    };
    
    eventSource.onerror = function(event) {
        console.error('Error in event source:', event);
    };
}

// Poll execution status
async function pollExecutionStatus(executionId) {
    const pollInterval = setInterval(async () => {
        try {
            const response = await fetch('/api/execution/' + executionId);
            const execution = await response.json();
            
            updateResultsSummary(execution);
            updateResults(execution.results);
            updateCharts(execution.results);
            
            if (execution.status === 'completed') {
                clearInterval(pollInterval);
                if (eventSource) {
                    eventSource.close();
                    eventSource = null;
                }
                currentExecution = null;
                updateActiveExecutions();
            }
            
        } catch (error) {
            console.error('Error polling execution status:', error);
            clearInterval(pollInterval);
        }
    }, 1000);
}

// Add log entry
function addLogEntry(level, message) {
    const logsContainer = document.getElementById('logs-container');
    const placeholder = logsContainer.querySelector('.log-placeholder');
    
    if (placeholder) {
        placeholder.remove();
    }
    
    const logEntry = document.createElement('div');
    logEntry.className = 'log-entry log-' + level;
    
    const timestamp = new Date().toLocaleTimeString();
    logEntry.innerHTML = '<span style="opacity: 0.7">[' + timestamp + ']</span> ' + message;
    
    logsContainer.appendChild(logEntry);
    logsContainer.scrollTop = logsContainer.scrollHeight;
}

// Clear logs
function clearLogs() {
    const logsContainer = document.getElementById('logs-container');
    logsContainer.innerHTML = '<div class="log-placeholder"><p>👋 Logs aparecerão aqui durante a execução</p></div>';
}

// Update results summary
function updateResultsSummary(execution) {
    const summaryContainer = document.getElementById('results-summary');
    summaryContainer.innerHTML = 
        '<div class="summary-item">' +
            '<span class="summary-number" style="color: #55efc4">' + execution.success_count + '</span>' +
            '<span class="summary-label">Sucessos</span>' +
        '</div>' +
        '<div class="summary-item">' +
            '<span class="summary-number" style="color: #ff7675">' + execution.failure_count + '</span>' +
            '<span class="summary-label">Falhas</span>' +
        '</div>' +
        '<div class="summary-item">' +
            '<span class="summary-number" style="color: #4facfe">' + execution.total_tests + '</span>' +
            '<span class="summary-label">Total</span>' +
        '</div>' +
        '<div class="summary-item">' +
            '<span class="summary-number" style="color: #fdcb6e">' + execution.status + '</span>' +
            '<span class="summary-label">Status</span>' +
        '</div>';
}

// Update results
function updateResults(results) {
    const resultsContainer = document.getElementById('results-container');
    const placeholder = resultsContainer.querySelector('.results-placeholder');
    
    if (placeholder) {
        placeholder.remove();
    }
    
    resultsContainer.innerHTML = '';
    
    results.forEach(result => {
        addResultItem(result);
    });
}

// Add single result
function addSingleResult(result) {
    const resultsContainer = document.getElementById('results-container');
    const placeholder = resultsContainer.querySelector('.results-placeholder');
    
    if (placeholder) {
        placeholder.remove();
    }
    
    addResultItem(result);
    showTab('results');
}

// Add result item
function addResultItem(result) {
    const resultsContainer = document.getElementById('results-container');
    
    const resultItem = document.createElement('div');
    resultItem.className = 'result-item';
    
    const statusClass = 'status-' + result.status;
    const httpStatus = result.http_status ? ' (HTTP ' + result.http_status + ')' : '';
    const error = result.error ? '<div style="color: #ff7675; margin-top: 5px;">Erro: ' + result.error + '</div>' : '';
    
    resultItem.innerHTML = 
        '<div class="result-header">' +
            '<strong>' + result.test_name + '</strong>' +
            '<div>' +
                '<span class="result-status ' + statusClass + '">' + result.status.toUpperCase() + httpStatus + '</span>' +
                '<span style="margin-left: 10px; color: #a0a6b8;">' + result.duration_ms + 'ms</span>' +
            '</div>' +
        '</div>' +
        '<div class="result-details">' +
            '<div>Executado em: ' + new Date(result.timestamp).toLocaleString() + '</div>' +
            error +
            (result.response ? '<details style="margin-top: 10px;"><summary>Resposta</summary><pre style="background: rgba(0,0,0,0.3); padding: 10px; border-radius: 4px; margin-top: 5px; overflow-x: auto;">' + result.response + '</pre></details>' : '') +
        '</div>';
    
    resultsContainer.appendChild(resultItem);
}

// Clear results
function clearResults() {
    const resultsContainer = document.getElementById('results-container');
    resultsContainer.innerHTML = '<div class="results-placeholder"><p>📊 Os resultados aparecerão aqui após a execução</p></div>';
    
    const summaryContainer = document.getElementById('results-summary');
    summaryContainer.innerHTML = '';
}

// Update active executions counter
function updateActiveExecutions() {
    const counter = currentExecution ? 1 : 0;
    document.getElementById('active-executions').textContent = counter;
}

// Initialize charts
function initializeCharts() {
    // Charts placeholder - simple text-based display
}

// Update charts
function updateCharts(results) {
    if (!results || results.length === 0) return;
    
    // Simple chart updates
    const statusCanvas = document.getElementById('status-chart');
    const timingCanvas = document.getElementById('timing-chart');
    
    if (statusCanvas) {
        const ctx = statusCanvas.getContext('2d');
        ctx.clearRect(0, 0, statusCanvas.width, statusCanvas.height);
        
        const successCount = results.filter(r => r.status === 'success').length;
        const failureCount = results.filter(r => r.status === 'failure').length;
        const errorCount = results.filter(r => r.status === 'error').length;
        
        ctx.fillStyle = '#e0e6ed';
        ctx.font = '14px Arial';
        ctx.textAlign = 'center';
        ctx.fillText('Sucessos: ' + successCount, statusCanvas.width/2, statusCanvas.height/2 - 20);
        ctx.fillText('Falhas: ' + failureCount, statusCanvas.width/2, statusCanvas.height/2);
        ctx.fillText('Erros: ' + errorCount, statusCanvas.width/2, statusCanvas.height/2 + 20);
    }
    
    if (timingCanvas) {
        const ctx = timingCanvas.getContext('2d');
        ctx.clearRect(0, 0, timingCanvas.width, timingCanvas.height);
        
        if (results.length > 0) {
            const maxValue = Math.max(...results.map(r => r.duration_ms));
            ctx.fillStyle = '#e0e6ed';
            ctx.font = '14px Arial';
            ctx.textAlign = 'center';
            ctx.fillText('Máx: ' + maxValue + 'ms', timingCanvas.width/2, timingCanvas.height/2 - 10);
            ctx.fillText('Testes: ' + results.length, timingCanvas.width/2, timingCanvas.height/2 + 10);
        }
    }
}

// Cleanup on page unload
window.addEventListener('beforeunload', function() {
    if (eventSource) {
        eventSource.close();
    }
});

// Drag and Drop Functions
function initializeDragAndDrop() {
    const testList = document.getElementById('test-list');
    if (!testList) return;
    
    // Add drag event listeners to test items
    const testItems = testList.querySelectorAll('.test-item');
    testItems.forEach(item => {
        item.addEventListener('dragstart', handleDragStart);
        item.addEventListener('dragover', handleDragOver);
        item.addEventListener('drop', handleDrop);
        item.addEventListener('dragend', handleDragEnd);
        item.addEventListener('dragenter', handleDragEnter);
        item.addEventListener('dragleave', handleDragLeave);
    });
}

function handleDragStart(e) {
    draggedElement = e.target;
    e.target.classList.add('dragging');
    e.dataTransfer.effectAllowed = 'move';
    e.dataTransfer.setData('text/html', e.target.outerHTML);
}

function handleDragOver(e) {
    if (e.preventDefault) {
        e.preventDefault();
    }
    e.dataTransfer.dropEffect = 'move';
    return false;
}

function handleDragEnter(e) {
    if (e.target.classList.contains('test-item') && e.target !== draggedElement) {
        e.target.classList.add('drag-over');
    }
}

function handleDragLeave(e) {
    if (e.target.classList.contains('test-item')) {
        e.target.classList.remove('drag-over');
    }
}

function handleDrop(e) {
    if (e.stopPropagation) {
        e.stopPropagation();
    }
    
    const targetElement = e.target.closest('.test-item');
    if (targetElement && targetElement !== draggedElement) {
        const testList = document.getElementById('test-list');
        const draggedIndex = Array.from(testList.children).indexOf(draggedElement);
        const targetIndex = Array.from(testList.children).indexOf(targetElement);
        
        if (draggedIndex < targetIndex) {
            targetElement.parentNode.insertBefore(draggedElement, targetElement.nextSibling);
        } else {
            targetElement.parentNode.insertBefore(draggedElement, targetElement);
        }
        
        updateExecutionOrder();
    }
    
    return false;
}

function handleDragEnd(e) {
    e.target.classList.remove('dragging');
    
    // Remove drag-over class from all items
    document.querySelectorAll('.test-item').forEach(item => {
        item.classList.remove('drag-over');
    });
    
    draggedElement = null;
}

// Execution Order Functions
function updateExecutionOrder() {
    const checkedBoxes = document.querySelectorAll('.test-checkbox:checked');
    const testList = document.getElementById('test-list');
    const allItems = Array.from(testList.querySelectorAll('.test-item'));
    
    // Build execution order based on DOM order and checked status
    executionOrder = [];
    allItems.forEach(item => {
        const checkbox = item.querySelector('.test-checkbox');
        if (checkbox && checkbox.checked) {
            executionOrder.push(checkbox.value);
        }
    });
    
    // Update counter
    const counter = document.getElementById('selected-count');
    if (counter) {
        counter.textContent = executionOrder.length;
        
        // Add visual feedback for counter
        if (executionOrder.length > 0) {
            counter.style.background = '#4facfe';
            counter.style.animation = 'pulse 0.5s ease-in-out';
        } else {
            counter.style.background = '#666';
            counter.style.animation = 'none';
        }
    }
    
    // Add visual indicators for execution order
    updateExecutionIndicators();
    
    // Log order change if there are selected tests
    if (executionOrder.length > 0) {
        console.log('🔄 Ordem de execução atualizada:', executionOrder);
    }
}

function updateExecutionIndicators() {
    const testItems = document.querySelectorAll('.test-item');
    testItems.forEach((item, index) => {
        const checkbox = item.querySelector('.test-checkbox');
        const testName = item.querySelector('.test-name');
        
        // Remove existing order indicator
        const existingIndicator = item.querySelector('.execution-indicator');
        if (existingIndicator) {
            existingIndicator.remove();
        }
        
        if (checkbox && checkbox.checked) {
            const orderIndex = executionOrder.indexOf(checkbox.value);
            if (orderIndex !== -1) {
                const indicator = document.createElement('span');
                indicator.className = 'execution-indicator';
                indicator.textContent = orderIndex + 1;
                indicator.style.cssText = `
                    background: #4facfe;
                    color: white;
                    border-radius: 50%;
                    width: 20px;
                    height: 20px;
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    font-size: 0.7rem;
                    font-weight: bold;
                    margin-left: auto;
                `;
                item.querySelector('.test-header').appendChild(indicator);
            }
        }
    });
}

function resetExecutionOrder() {
    // Uncheck all checkboxes
    document.querySelectorAll('.test-checkbox').forEach(checkbox => {
        checkbox.checked = false;
    });
    
    // Reset execution order array
    executionOrder = [];
    
    // Update UI
    updateExecutionOrder();
    
    // Show feedback
    addLogEntry('info', '🔄 Ordem de execução resetada');
    
    // Visual feedback on reset button
    const resetBtn = event.target;
    const originalText = resetBtn.textContent;
    resetBtn.textContent = '✓';
    resetBtn.style.background = '#27ae60';
    setTimeout(() => {
        resetBtn.textContent = originalText;
        resetBtn.style.background = '';
    }, 1000);
}

function selectAllTests() {
    document.querySelectorAll('.test-checkbox').forEach(checkbox => {
        checkbox.checked = true;
    });
    updateExecutionOrder();
}

function clearSelection() {
    document.querySelectorAll('.test-checkbox').forEach(checkbox => {
        checkbox.checked = false;
    });
    updateExecutionOrder();
}