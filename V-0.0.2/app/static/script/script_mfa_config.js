function validateMFAToken() {

    fetch('/validate_mfa_token', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ token: document.getElementById('mfa_token').value })
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