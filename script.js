document.getElementById('submit').addEventListener('click', function() {
    const userprompt = document.getElementById('prompt').value;
    const peterInfo = document.querySelector('.peter-info');
    
    if (!userprompt.trim()) {
        document.getElementById('response-container').innerHTML = `
            <div class="response-item">
                <div class="title">Ehehehehe...</div>
                <div class="content">Hey, you gotta type something first! Like that time I tried to send an empty letter to Santa...</div>
            </div>`;
        return;
    }

    // Add loading state
    peterInfo.classList.add('loading');
    document.getElementById('response-container').innerHTML = `
        <div class="response-item">
            <div class="title">Hold on a sec...</div>
            <div class="content">Peter is thinking (which usually takes a while, hehehe)...</div>
        </div>`;

    fetch('http://localhost:8080/prompt', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            prompt: userprompt
        })
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        return response.json();
    })
    .then(data => {
        peterInfo.classList.remove('loading');
        if (data.result) {
            const responseHTML = `
                <div class="response-item">
                    <div class="title">Peter Says:</div>
                    <div class="content">${data.result}</div>
                </div>`;
            document.getElementById('response-container').innerHTML = responseHTML;
        } else {
            throw new Error('No response received');
        }
    })
    .catch(err => {
        peterInfo.classList.remove('loading');
        console.error("Error:", err);
        const errorHTML = `
            <div class="response-item">
                <div class="title">Oh crap!</div>
                <div class="content">Something went wrong! Like that time I tried to fix the TV with a hammer... Maybe try again?</div>
            </div>`;
        document.getElementById('response-container').innerHTML = errorHTML;
    });
});

// Clear response when starting to type new message
document.getElementById('prompt').addEventListener('focus', function() {
    if (document.querySelector('.peter-info').classList.contains('loading')) {
        return; // Don't clear if still loading
    }
    document.getElementById('response-container').innerHTML = '';
});