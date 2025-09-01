function validateMFAToken() {

    fetch('/first_validate_mfa', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },

        // send form data
        body: JSON.stringify({
            username: document.getElementById('username').value,
            mfaSecret: document.getElementById('mfaSecret').value,
            mfaCode: document.getElementById('mfaCode').value
        })
        
    })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                // print success message and redirect to home
                alert('MFA token validated successfully!');
                window.location.href = '/';  // Redirect to home on success';
            } else {
                document.getElementById('error_message').textContent = 'Invalid token. Please try again.';
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
    return false; // Prevent form submission

}