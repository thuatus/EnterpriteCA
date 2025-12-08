// This script is for the search functionality
document.addEventListener('DOMContentLoaded', function () {
    // load certificates when the page loads
    loadCertificates();

});

const searchInput = document.getElementById('searchInput');
const table = document.getElementById('certTable').getElementsByTagName('tbody')[0];

searchInput.focus();


document.querySelector("#certTable tbody").addEventListener("click", function(event) {
  const target = event.target;

  // Botão "View Certificate"
  if (target.classList.contains("view-cert-btn")) {
    const modalId = target.dataset.modalId;
    const overlayId = target.dataset.overlayId;
    document.getElementById(modalId).style.display = "block";
    document.getElementById(overlayId).style.display = "block";
  }

  // Botão "Revoke Certificate"
  if (target.classList.contains("revoke-btn")) {
    const subject = target.dataset.subject;
    revokeCert(subject); // Certifique-se de que essa função está definida
  }

  // Botão "View Key"
  if (target.classList.contains("btn-primary")) {
    viewKey(target); // Repassa o botão como referência
  }

  // Botão "Copy Certificate"
  if (target.textContent.includes("Copy Certificate")) {
    const pre = target.previousElementSibling;
    if (pre && pre.tagName === "PRE") {
      copyCert(pre.id);
    }
  }

  // Botão "Copy Key"
  if (target.textContent.includes("Copy Key")) {
    copyKey(); // Supondo que essa função já esteja definida
  }

  // Botão de fechar modal (X)
    if (target.classList.contains("close-key-modal")) {
        document.getElementById('modalViewKey').style.display = 'none';
    }       

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
                        <button class="btn view-cert-btn" data-modal-id="modal-${cert.subject}" data-overlay-id="overlay-${cert.subject}">View Certificate </button>
                        <div class="overlay" id="overlay-${cert.subject}" onclick="closeModal('modal-${cert.subject}', 'overlay-${cert.subject}')"></div>
                        <div class="modal" id="modal-${cert.subject}">
                            <p><strong>Public Key:</strong></p>
                            <pre id="publicKeyText-${cert.subject}">${cert.public_key}</pre>
                            <button class="btn" onclick="copyCert('publicKeyText-${cert.subject}')">Copy Certificate</button>
                        </div>
                        <button class="btn btn-primary" onclick="viewKey(this)">View Key</button>
                        <div id="modalViewKey" class="modal">
                            <div class="modal-content">
                                <span class="close-key-modal">&times;</span>
                                <p id="key_content"></p>
                                <button class="btn" onclick="copyKey()">Copy Key</button>
                            </div>
                        </div>

                        <button class="btn-red revoke-btn" data-subject="${cert.subject}" onclick="revokeCert(${cert.subject})">Revoke Certificate</button>
                    </td>
                `;
                tbody.appendChild(row);
            });
        })
        .catch(error => {
            console.error('Error fetching certificates:', error);
        });



}

searchInput.addEventListener('keyup', function () {
    const filter = searchInput.value.toLowerCase();
    const rows = table.getElementsByTagName('tr');
    for (let i = 0; i < rows.length; i++) {
        const rowText = rows[i].textContent.toLowerCase();
        rows[i].classList.toggle('hide', !rowText.includes(filter));
    }
});

function showModal(modalId, overlayId) {
    document.getElementById(modalId).style.display = 'block';
    document.getElementById(overlayId).style.display = 'block';
}

function closeModal(modalId, overlayId) {
    document.getElementById(modalId).style.display = 'none';
    document.getElementById(overlayId).style.display = 'none';
}

function copyCert(publicKeyTextId) {
    const text = document.getElementById(publicKeyTextId).innerText;
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



// Copy key content to clipboard
function copyKey() {
    const text = document.getElementById('key_content').innerText;
    navigator.clipboard.writeText(text).then(() => {
        alert('Content copied to clipboard!');
    });
}

// Send Request to revoke certificate
function revokeCert(subject) {
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
}

