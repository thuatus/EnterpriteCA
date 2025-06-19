    
        function confirmSubmit() {
            // retund aks with site name confirmation
            var serverName = document.getElementById('serverName').value;
            if (!serverName) {
                alert('Please enter a server name.');
                return false;
            }
            // Confirm with the user
            if (serverName.length < 3) {
                alert('Server name must be at least 3 characters long.');
                return false;
            }
            if (!/^[a-zA-Z0-9.-]+$/.test(serverName)) {
                alert('Server name can only contain letters, numbers, dots, and hyphens.');
                return false;
            }
            if (!serverName.endsWith('.com') &&
                !serverName.endsWith('.org') &&
                !serverName.endsWith('.net')) {
                alert('Server name must end with .com, .org, or .net.');
                return false;
            }
            // If all checks pass, confirm submission
            if (confirm('Are you sure you want to issue a certificate for ' + serverName + '?')) {
                return true;
            }
            // If user cancels, prevent form submission
            return false;
        }

    