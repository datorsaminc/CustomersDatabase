// API communication layer - handles all HTTP requests to the backend
const Api = {
    baseUrl: '/api/customers',

    // Fetch customers with optional search, page and limit parameters
    async getCustomers(search = '', page = 1, limit = 50) {
        const params = new URLSearchParams({ page, limit });
        if (search) params.append('search', search);
        
        try {
            const response = await fetch(`${this.baseUrl}?${params.toString()}`);
            if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
            return await response.json();
        } catch (error) {
            console.error('Error fetching customers:', error);
            throw error;
        }
    },

    // Fetch a single customer by ID
    async getCustomer(id) {
        try {
            const response = await fetch(`${this.baseUrl}/${id}`);
            if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
            return await response.json();
        } catch (error) {
            console.error('Error fetching customer:', error);
            throw error;
        }
    },

    // Create a new customer
    async createCustomer(customerData) {
        try {
            const response = await fetch(this.baseUrl, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(customerData)
            });
            if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
            return await response.json();
        } catch (error) {
            console.error('Error creating customer:', error);
            throw error;
        }
    },

    // Update an existing customer
    async updateCustomer(id, customerData) {
        try {
            const response = await fetch(`${this.baseUrl}/${id}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(customerData)
            });
            if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
            return await response.json();
        } catch (error) {
            console.error('Error updating customer:', error);
            throw error;
        }
    },

    // Delete a customer by ID
    async deleteCustomer(id) {
    	try {
    		const response = await fetch(`${this.baseUrl}/${id}`, {
    			method: 'DELETE'
    		});
    		if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
    		return await response.json();
    	} catch (error) {
    		console.error('Error deleting customer:', error);
    		throw error;
    	}
    },
   
    // Get total licenses sold statistics
    async getLicenseStats() {
    	try {
    		const response = await fetch('/api/licenses/stats');
    		if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
    		return await response.json();
    	} catch (error) {
    		console.error('Error fetching license stats:', error);
    		throw error;
    	}
    }
   };
