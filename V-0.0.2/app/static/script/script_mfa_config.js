function validateMFAToken() {
    fetch('/first_validate_mfa', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            username: document.getElementById('username').value,
            mfaSecret: document.getElementById('mfaSecret').value,
            mfaCode: document.getElementById('mfaCode').value
        })
    })
    .then(response => {
        if (!response.ok) throw new Error('Network response was not ok');
        return response.json();
    })
    .then(data => {
        if (data.success) {
            alert('MFA token validated successfully!');
            window.location.href = '/';
        } else {
            document.getElementById('error_message').textContent = 'Invalid token. Please try again.';
        }
    })
    .catch(error => {
        document.getElementById('error_message').textContent = 'Error validating MFA token.';
        console.error('Error:', error);
    });
    return false; // Prevent form submission
}