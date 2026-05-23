// Main application entry point
$(document).ready(function() {
    // Initialize all components
    Search.init();
    CustomerForm.init();

    // Load initial customer list and license stats
    CustomerTable.render('');
    loadLicenseStats();

    // Event delegation for edit buttons (dynamically added)
    $(document).on('click', '.edit-btn', function(e) {
        e.preventDefault();
        const id = $(this).data('id');
        CustomerForm.editCustomer(id);
    });

    // Event delegation for delete buttons
    $(document).on('click', '.delete-btn', function(e) {
        e.preventDefault();
        const id = $(this).data('id');
        
        // Get customer info from the row
        const $row = $(this).closest('tr');
        const company = $row.find('td:first').text().trim();
        const contactName = $row.find('td:nth-child(2)').text().trim();
        
        $('#deleteCustomerInfo').html(`<strong>${company}</strong><br>Contact: ${contactName}`);
        $('#confirmDeleteBtn').data('id', id);
        
        const deleteModal = new bootstrap.Modal(document.getElementById('deleteModal'));
        deleteModal.show();
    });

    // Confirm delete action
    $('#confirmDeleteBtn').on('click', async function() {
        const id = $(this).data('id');
        if (!id) return;

        try {
            await Api.deleteCustomer(id);
            
            // Hide the modal
            const deleteModalEl = document.getElementById('deleteModal');
            const deleteModalInstance = bootstrap.Modal.getInstance(deleteModalEl);
            deleteModalInstance.hide();
            
            CustomerTable.render(CustomerTable.currentSearch);
            loadLicenseStats();
            showToast('Customer deleted successfully', 'success');
        } catch (error) {
            showToast('Failed to delete customer', 'danger');
        }
    });

    // Event delegation for pagination clicks
    $(document).on('click', '#paginationControls .page-link', function(e) {
        e.preventDefault();
        if ($(this).parent().hasClass('disabled')) return;
        
        const page = parseInt($(this).data('page'));
        CustomerTable.onPageClick(page);
    });
});

// Show a toast notification
function showToast(message, type = 'info') {
    const $toastEl = $('#notificationToast');
    
    // Remove any existing background color classes
    $toastEl.removeClass('bg-success bg-danger bg-warning bg-info bg-primary');
    
    // Add appropriate background class based on type
    switch (type) {
        case 'success':
            $toastEl.addClass('bg-success');
            break;
        case 'danger':
            $toastEl.addClass('bg-danger');
            break;
        case 'warning':
            $toastEl.addClass('bg-warning text-dark');
            break;
        default:
            $toastEl.addClass('bg-info');
    }

    $('#toastMessage').text(message);
    
    const toast = new bootstrap.Toast($toastEl, { delay: 3000 });
    toast.show();
}

// Load and display total licenses sold
async function loadLicenseStats() {
    try {
        const stats = await Api.getLicenseStats();
        $('#totalLicensesSold').text(stats.totalLicenses.toLocaleString());
    } catch (error) {
        console.error('Failed to load license stats:', error);
        $('#totalLicensesSold').text('-');
    }
}
