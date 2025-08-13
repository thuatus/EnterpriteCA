// This script is for the search functionality
document.addEventListener('DOMContentLoaded', function () {
        // load certificates when the page loads
    loadCertificates();

    const searchInput = document.getElementById('searchInput');
    const table = document.getElementById('certTable').getElementsByTagName('tbody')[0];

    searchInput.focus();

    searchInput.addEventListener('keyup', function () {
        const filter = searchInput.value.toLowerCase();
        const rows = table.getElementsByTagName('tr');
        for (let i = 0; i < rows.length; i++) {
            const rowText = rows[i].textContent.toLowerCase();
            rows[i].classList.toggle('hide', !rowText.includes(filter));
        }
    });


    document.querySelectorAll('.view-cert-btn').forEach(function (btn) {
        btn.addEventListener('click', function () {
            document.getElementById('publicKeyText').innerText = btn.getAttribute('data-public-key');
            document.getElementById('modal').style.display = 'block';
            document.getElementById('overlay').style.display = 'block';
        });
    });

    function loadCertificates() {
        fetch('/view_cert_info/') // Adjust the endpoint if needed
            .then(response => response.json())
            .then(data => {
                const certs = data.certificates;
                const tbody = document.querySelector("#certTable tbody");
                tbody.innerHTML = ""; // Clear existing rows

                certs.forEach(cert => {
                    const row = document.createElement("tr");
                    row.innerHTML = `
                    <td id="subject_cell">${cert.subject}</td>
                    <td>${cert.created}</td>
                    <td>${cert.expire}</td>
                    <td>${cert.status}</td>
                    <td>
                        <button class="btn view-cert-btn" data-public-key="${cert.public_key}">View Certificate</button>
                        <div class="overlay" id="overlay" onclick="closeModal()"></div>
                        <div class="modal" id="modal">
                            <p><strong>Public Key:</strong></p>
                            <pre id="publicKeyText"></pre>
                            <button class="btn" onclick="copyCert()">Copy Certificate</button>
                        </div>
                        <button class="btn btn-primary" onclick="viewKey(this)">View Key</button>
                        <div id="modalViewKey" class="modal">
                            <div class="modal-content">
                                <span class="close">&times;</span>
                                <p id="key_content"></p>
                                <button class="btn" onclick="copyKey()">Copy Key</button>
                            </div>
                        </div>

                        <button class="btn-red revoke-btn" data-subject="${cert.Subject}">Revoke Certificate</button>
                    </td>
                `;
                    tbody.appendChild(row);
                });
            })
            .catch(error => {
                console.error('Error fetching certificates:', error);
            });
    }




    function showModal() {
        document.getElementById('modal').style.display = 'block';
        document.getElementById('overlay').style.display = 'block';
    }

    function closeModal() {
        document.getElementById('modal').style.display = 'none';
        document.getElementById('overlay').style.display = 'none';
    }

    function copyCert() {
        const text = document.getElementById('publicKeyText').innerText;
        navigator.clipboard.writeText(text).then(() => {
            alert('Content copied to clipboard!');
        });
    }
    function viewKey(btn) {

        // Get the subject from the same row as the clicked button
        const subject = btn.closest('tr').querySelector('#subject_cell').innerText;

        fetch('/view_server_key/', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ subject })
        })
            .then(response => response.json())
            .then(data => {
                document.getElementById('key_content').innerText = data.Server_key || 'No key found';
                document.getElementById('modalViewKey').style.display = 'block';
            })
            .catch(error => {
                document.getElementById('key_content').innerText = 'Error fetching key';
                document.getElementById('modalViewKey').style.display = 'block';
            });
    }

    // Close modal when clicking the close button
    document.addEventListener('DOMContentLoaded', function () {
        const closeBtn = document.querySelector('#modalViewKey .close');
        if (closeBtn) {
            closeBtn.onclick = function () {
                document.getElementById('modalViewKey').style.display = 'none';
            };
        }
    });

    // Copy key content to clipboard
    function copyKey() {
        const text = document.getElementById('key_content').innerText;
        navigator.clipboard.writeText(text).then(() => {
            alert('Content copied to clipboard!');
        });
    }

    // Send Request to revoke certificate
    document.querySelectorAll('.revoke-btn').forEach(function (btn) {
        btn.addEventListener('click', function () {
            // include a confirmation dialog

            const subject = btn.getAttribute('data-subject');

            if (!confirm(`Are you sure you want to revoke this certificate: ${subject}?`)) {
                return; // If the user cancels, do nothing
            }
            fetch('/revoke_cert/', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ subject })
            })
                .then(response => {
                    if (response.ok) {
                        alert('Success on revoke certificate');
                    } else {
                        alert('Failed to revoke certificate');
                    }
                })
                .catch(() => alert('Error communicating with server'));
        });
    });


});