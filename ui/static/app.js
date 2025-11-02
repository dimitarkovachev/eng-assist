// State management
let isPaused = false;

// DOM elements
const transcriptElement = document.getElementById('transcript');
const responseElement = document.getElementById('response');
const costElement = document.querySelector('.cost');

// Keyboard shortcuts
document.addEventListener('keydown', (e) => {
    if (e.ctrlKey) {
        switch(e.key.toLowerCase()) {
            case 'r':
                e.preventDefault();
                resetAssistant();
                break;
            case 'p':
                e.preventDefault();
                togglePause();
                break;
            case 'q':
                e.preventDefault();
                // Quit functionality can be handled by closing the tab
                break;
        }
    }
});

// UI Actions
function resetAssistant() {
    fetch('/api/reset', { method: 'POST' })
        .then(() => {
            transcriptElement.textContent = '';
            responseElement.textContent = '';
        })
        .catch(console.error);
}

function togglePause() {
    isPaused = !isPaused;
    fetch('/api/pause', { method: 'POST' })
        .catch(console.error);
}

// State polling
function updateState() {
    fetch('/api/state')
        .then(response => response.json())
        .then(data => {
            if (data.transcript !== undefined) {
                transcriptElement.textContent = data.transcript;
            }
            if (data.response !== undefined) {
                responseElement.textContent = data.response;
            }
            if (data.cost !== undefined) {
                costElement.textContent = `Cost: $${data.cost.toFixed(4)}`;
            }
        })
        .catch(console.error);
}

// Start polling
setInterval(updateState, 250); 