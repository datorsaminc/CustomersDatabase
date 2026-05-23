// Customer form handling with Bootstrap validation
const CustomerForm = {
    modal: null,
    isEditing: false,
    currentId: null,

    init() {
        this.modal = new bootstrap.Modal(document.getElementById('customerModal'));
        
        // Form submission handler
        $('#customerForm').on('submit', (e) => {
            e.preventDefault();
            
            if (!this.validateForm()) return;
            
            const customerData = this.getFormData();
            
            if (this.isEditing && this.currentId) {
                this.updateCustomer(customerData);
            } else {
                this.createCustomer(customerData);
            }
        });

        // Reset form when modal is closed
        $('#customerModal').on('hidden.bs.modal', () => {
            this.resetForm();
        });

        // Add customer button
        $('#addCustomerBtn').on('click', () => {
            this.isEditing = false;
            this.currentId = null;
            $('#customerModalLabel').text('Add New Customer');
            this.resetForm();
            this.modal.show();
        });
    },

    // Open form for editing an existing customer
    async editCustomer(id) {
        try {
            const customer = await Api.getCustomer(id);
            
            this.isEditing = true;
            this.currentId = id;
            $('#customerModalLabel').text('Edit Customer');
            
            // Populate form fields
            $('#customerId').val(customer.id);
            $('#company').val(customer.company || '');
            $('#name1').val(customer.name1 || '');
            $('#name2').val(customer.name2 || '');
            $('#email').val(customer.email || '');
            $('#programVersion').val(customer.programVersion || '');
            $('#deliveryDate').val(customer.deliveryDate || '');
            $('#licenses').val(customer.licenses != null ? customer.licenses : '');
            $('#visitAddress').val(customer.visitAddress || '');
            $('#mailingAddress').val(customer.mailingAddress || '');
            $('#postalCodeCity').val(customer.postalCodeCity || '');
            $('#landlinePhone').val(customer.landlinePhone || '');
            $('#mobilePhone').val(customer.mobilePhone || '');
            $('#faxNumber').val(customer.faxNumber || '');
            $('#comments').val(customer.comments || '');
            
            // Remove validation errors
            $('#customerForm').removeClass('was-validated');
            
            this.modal.show();
        } catch (error) {
            showToast('Failed to load customer details', 'danger');
        }
    },

    // Create a new customer
    async createCustomer(customerData) {
        try {
            await Api.createCustomer(customerData);
            this.modal.hide();
            CustomerTable.render(CustomerTable.currentSearch);
            loadLicenseStats();
            showToast('Customer created successfully', 'success');
        } catch (error) {
            showToast('Failed to create customer', 'danger');
        }
    },

    // Update an existing customer
    async updateCustomer(customerData) {
        try {
            await Api.updateCustomer(this.currentId, customerData);
            this.modal.hide();
            CustomerTable.render(CustomerTable.currentSearch);
            loadLicenseStats();
            showToast('Customer updated successfully', 'success');
        } catch (error) {
            showToast('Failed to update customer', 'danger');
        }
    },

    // Get form data as an object
    getFormData() {
        const licensesVal = parseInt($('#licenses').val(), 10);
        return {
            company: $('#company').val().trim(),
            name1: $('#name1').val().trim(),
            name2: $('#name2').val().trim() || null,
            email: $('#email').val().trim(),
            programVersion: $('#programVersion').val().trim(),
            deliveryDate: $('#deliveryDate').val().trim(),
            licenses: isNaN(licensesVal) || licensesVal <= 0 ? null : licensesVal,
            visitAddress: $('#visitAddress').val().trim(),
            mailingAddress: $('#mailingAddress').val().trim(),
            postalCodeCity: $('#postalCodeCity').val().trim(),
            landlinePhone: $('#landlinePhone').val().trim(),
            mobilePhone: $('#mobilePhone').val().trim() || null,
            faxNumber: $('#faxNumber').val().trim() || null,
            comments: $('#comments').val().trim()
        };
    },

    // Validate form fields using Bootstrap's validation classes
    validateForm() {
        const $form = $('#customerForm');
        let isValid = true;

        // Company - required, min 2 characters
        const company = $('#company').val().trim();
        if (company.length < 2) {
            $('#company').addClass('is-invalid').removeClass('is-valid');
            isValid = false;
        } else {
            $('#company').addClass('is-valid').removeClass('is-invalid');
        }

        // Name1 - required, min 2 characters
        const name1 = $('#name1').val().trim();
        if (name1.length < 2) {
            $('#name1').addClass('is-invalid').removeClass('is-valid');
            isValid = false;
        } else {
            $('#name1').addClass('is-valid').removeClass('is-invalid');
        }

        // Email - optional, but must be valid format if provided
        const email = $('#email').val().trim();
        if (email) {
            const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            if (!emailRegex.test(email)) {
                $('#email').addClass('is-invalid').removeClass('is-valid');
                isValid = false;
            } else {
                $('#email').addClass('is-valid').removeClass('is-invalid');
            }
        } else {
            $('#email').removeClass('is-invalid is-valid');
        }

        return isValid;
    },

    // Reset form to initial state
    resetForm() {
        $('#customerForm')[0].reset();
        $('#customerId').val('');
        $('#customerForm').removeClass('was-validated');
        $('.form-control').removeClass('is-valid is-invalid');
        this.isEditing = false;
        this.currentId = null;
    }
};
