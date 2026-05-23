// Customer table rendering and interactions
const CustomerTable = {
    currentPage: 1,
    limit: 50,
    currentSearch: '',

    // Render the customer table with pagination
    async render(search = '') {
        this.currentSearch = search;
        try {
            const data = await Api.getCustomers(search, this.currentPage, this.limit);
            
            $('#showingCount').text(data.customers.length);
            $('#totalCount').text(data.total);

            const tbody = $('#customerTableBody');
            tbody.empty();

            if (data.customers.length === 0) {
                tbody.append(`
                    <tr>
                        <td colspan="7" class="text-center text-muted py-4">No customers found</td>
                    </tr>
                `);
            } else {
                data.customers.forEach(customer => {
                    const licenses = customer.licenses != null ? customer.licenses : '';
                    const row = $(`
                        <tr data-id="${customer.id}">
                            <td><strong>${this.escapeHtml(customer.company)}</strong></td>
                            <td>${this.escapeHtml(customer.name1)}${customer.name2 ? '<br><small class="text-muted">' + this.escapeHtml(customer.name2) + '</small>' : ''}</td>
                            <td>${this.escapeHtml(customer.email)}</td>
                            <td>${this.escapeHtml(customer.landlinePhone || '')}${customer.mobilePhone ? '<br><small class="text-muted">M: ' + this.escapeHtml(customer.mobilePhone) + '</small>' : ''}</td>
                            <td>${this.escapeHtml(customer.postalCodeCity)}</td>
                            <td class="text-center">${licenses}</td>
                            <td class="text-end">
                                <button class="btn btn-sm btn-outline-primary edit-btn" data-id="${customer.id}" title="Edit">&#9998;</button>
                                <button class="btn btn-sm btn-outline-danger delete-btn" data-id="${customer.id}" title="Delete">&#128465;</button>
                            </td>
                        </tr>
                    `);
                    tbody.append(row);
                });
            }

            this.renderPagination(data.total, data.page, this.limit);
        } catch (error) {
            showToast('Failed to load customers', 'danger');
        }
    },

    // Render pagination controls
    renderPagination(totalItems, currentPage, limit) {
        const totalPages = Math.ceil(totalItems / limit);
        const $pagination = $('#paginationControls');
        $pagination.empty();

        if (totalPages <= 1) return;

        // Previous button
        const prevDisabled = currentPage === 1 ? 'disabled' : '';
        $pagination.append(`
            <li class="page-item ${prevDisabled}">
                <a class="page-link" href="#" data-page="${currentPage - 1}">Previous</a>
            </li>
        `);

        // Page numbers (show limited range)
        const startPage = Math.max(1, currentPage - 2);
        const endPage = Math.min(totalPages, currentPage + 2);

        if (startPage > 1) {
            $pagination.append(`
                <li class="page-item">
                    <a class="page-link" href="#" data-page="1">1</a>
                </li>
            `);
            if (startPage > 2) {
                $pagination.append('<li class="page-item disabled"><span class="page-link">...</span></li>');
            }
        }

        for (let i = startPage; i <= endPage; i++) {
            const activeClass = i === currentPage ? 'active' : '';
            $pagination.append(`
                <li class="page-item ${activeClass}">
                    <a class="page-link" href="#" data-page="${i}">${i}</a>
                </li>
            `);
        }

        if (endPage < totalPages) {
            if (endPage < totalPages - 1) {
                $pagination.append('<li class="page-item disabled"><span class="page-link">...</span></li>');
            }
            $pagination.append(`
                <li class="page-item">
                    <a class="page-link" href="#" data-page="${totalPages}">${totalPages}</a>
                </li>
            `);
        }

        // Next button
        const nextDisabled = currentPage === totalPages ? 'disabled' : '';
        $pagination.append(`
            <li class="page-item ${nextDisabled}">
                <a class="page-link" href="#" data-page="${currentPage + 1}">Next</a>
            </li>
        `);
    },

    // Handle pagination clicks
    onPageClick(page) {
        this.currentPage = page;
        this.render(this.currentSearch);
    },

    // Escape HTML to prevent XSS
    escapeHtml(text) {
        if (!text) return '';
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
};
