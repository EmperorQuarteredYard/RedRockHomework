// 表格相关功能封装
class DataTable {
    constructor(tableId, options = {}) {
        this.table = document.getElementById(tableId);
        this.options = {
            searchable: true,
            pagination: true,
            pageSize: 10,
            sortable: true,
            ...options
        };

        this.data = [];
        this.filteredData = [];
        this.currentPage = 1;
        this.sortColumn = null;
        this.sortDirection = 'asc';

        this.init();
    }

    init() {
        if (!this.table) {
            console.error(`表格 ${this.tableId} 不存在`);
            return;
        }

        this.createHeader();
        this.createBody();

        if (this.options.searchable) {
            this.addSearch();
        }
    }

    createHeader() {
        const thead = this.table.querySelector('thead');
        if (!thead) return;

        const headerRow = thead.querySelector('tr');
        if (!headerRow) return;

        // 添加排序功能
        if (this.options.sortable) {
            const headers = headerRow.querySelectorAll('th[data-sort]');
            headers.forEach(th => {
                th.style.cursor = 'pointer';
                th.addEventListener('click', () => this.sortTable(th.dataset.sort));
            });
        }
    }

    createBody() {
        this.tbody = this.table.querySelector('tbody');
        if (!this.tbody) {
            this.tbody = document.createElement('tbody');
            this.table.appendChild(this.tbody);
        }
    }

    addSearch() {
        const container = this.table.closest('.table-container') || this.table.parentElement;

        const searchDiv = document.createElement('div');
        searchDiv.className = 'search-box';
        searchDiv.innerHTML = `
            <i class="search-icon">🔍</i>
            <input type="text" placeholder="搜索..." class="search-input">
        `;

        const searchInput = searchDiv.querySelector('.search-input');
        searchInput.addEventListener('input', debounce(() => {
            this.search(searchInput.value);
        }, 300));

        // 插入搜索框
        const toolbar = container.querySelector('.table-toolbar');
        if (toolbar) {
            toolbar.prepend(searchDiv);
        } else {
            container.insertBefore(searchDiv, this.table);
        }
    }

    setData(data) {
        this.data = data;
        this.filteredData = [...data];
        this.render();
    }

    render() {
        if (!this.tbody) return;

        this.tbody.innerHTML = '';

        const startIndex = (this.currentPage - 1) * this.options.pageSize;
        const endIndex = startIndex + this.options.pageSize;
        const pageData = this.filteredData.slice(startIndex, endIndex);

        pageData.forEach((row, index) => {
            const tr = document.createElement('tr');

            // 生成行数据
            const columns = this.getColumns();
            columns.forEach(column => {
                const td = document.createElement('td');

                if (typeof column === 'function') {
                    td.innerHTML = column(row, index);
                } else if (typeof column === 'string') {
                    td.textContent = this.getNestedValue(row, column);
                } else if (column.field) {
                    if (column.formatter) {
                        td.innerHTML = column.formatter(row[column.field], row, index);
                    } else {
                        td.textContent = this.getNestedValue(row, column.field);
                    }
                }

                tr.appendChild(td);
            });

            this.tbody.appendChild(tr);
        });

        if (this.options.pagination) {
            this.renderPagination();
        }
    }

    getColumns() {
        // 从表头获取列定义
        const thead = this.table.querySelector('thead');
        if (!thead) return [];

        const headers = thead.querySelectorAll('th');
        return Array.from(headers).map(th => ({
            field: th.dataset.field,
            title: th.textContent,
            width: th.style.width
        }));
    }

    getNestedValue(obj, path) {
        return path.split('.').reduce((o, p) => (o ? o[p] : ''), obj);
    }

    search(keyword) {
        if (!keyword.trim()) {
            this.filteredData = [...this.data];
        } else {
            this.filteredData = this.data.filter(row => {
                return Object.values(row).some(value =>
                    String(value).toLowerCase().includes(keyword.toLowerCase())
                );
            });
        }

        this.currentPage = 1;
        this.render();
    }

    sortTable(column) {
        if (this.sortColumn === column) {
            this.sortDirection = this.sortDirection === 'asc' ? 'desc' : 'asc';
        } else {
            this.sortColumn = column;
            this.sortDirection = 'asc';
        }

        this.filteredData.sort((a, b) => {
            const aValue = this.getNestedValue(a, column);
            const bValue = this.getNestedValue(b, column);

            if (aValue < bValue) return this.sortDirection === 'asc' ? -1 : 1;
            if (aValue > bValue) return this.sortDirection === 'asc' ? 1 : -1;
            return 0;
        });

        this.render();
    }

    renderPagination() {
        const totalPages = Math.ceil(this.filteredData.length / this.options.pageSize);
        if (totalPages <= 1) return;

        const paginationDiv = document.createElement('div');
        paginationDiv.className = 'pagination';

        // 上一页
        const prevBtn = document.createElement('button');
        prevBtn.className = `page-item ${this.currentPage === 1 ? 'disabled' : ''}`;
        prevBtn.textContent = '上一页';
        prevBtn.disabled = this.currentPage === 1;
        prevBtn.addEventListener('click', () => {
            if (this.currentPage > 1) {
                this.currentPage--;
                this.render();
            }
        });
        paginationDiv.appendChild(prevBtn);

        // 页码
        for (let i = 1; i <= totalPages; i++) {
            if (
                i === 1 ||
                i === totalPages ||
                (i >= this.currentPage - 1 && i <= this.currentPage + 1)
            ) {
                const pageBtn = document.createElement('button');
                pageBtn.className = `page-item ${i === this.currentPage ? 'active' : ''}`;
                pageBtn.textContent = i;
                pageBtn.addEventListener('click', () => {
                    this.currentPage = i;
                    this.render();
                });
                paginationDiv.appendChild(pageBtn);
            } else if (
                i === this.currentPage - 2 ||
                i === this.currentPage + 2
            ) {
                const ellipsis = document.createElement('span');
                ellipsis.textContent = '...';
                ellipsis.style.padding = '5px 10px';
                paginationDiv.appendChild(ellipsis);
            }
        }

        // 下一页
        const nextBtn = document.createElement('button');
        nextBtn.className = `page-item ${this.currentPage === totalPages ? 'disabled' : ''}`;
        nextBtn.textContent = '下一页';
        nextBtn.disabled = this.currentPage === totalPages;
        nextBtn.addEventListener('click', () => {
            if (this.currentPage < totalPages) {
                this.currentPage++;
                this.render();
            }
        });
        paginationDiv.appendChild(nextBtn);

        // 插入分页
        const container = this.table.closest('.table-container') || this.table.parentElement;
        const existingPagination = container.querySelector('.pagination');
        if (existingPagination) {
            existingPagination.replaceWith(paginationDiv);
        } else {
            container.appendChild(paginationDiv);
        }
    }

    updateRow(index, newData) {
        const dataIndex = (this.currentPage - 1) * this.options.pageSize + index;
        if (dataIndex >= 0 && dataIndex < this.filteredData.length) {
            this.filteredData[dataIndex] = { ...this.filteredData[dataIndex], ...newData };
            this.render();
        }
    }

    deleteRow(index) {
        const dataIndex = (this.currentPage - 1) * this.options.pageSize + index;
        if (dataIndex >= 0 && dataIndex < this.filteredData.length) {
            this.filteredData.splice(dataIndex, 1);
            this.data.splice(dataIndex, 1);
            this.render();
        }
    }
}

// 表格格式化函数
class TableFormatters {
    static date(value) {
        return UI.formatDate(value);
    }

    static status(value) {
        const statusMap = {
            active: { text: '激活', class: 'status-active' },
            inactive: { text: '停用', class: 'status-inactive' },
            pending: { text: '待审核', class: 'status-pending' }
        };

        const status = statusMap[value] || { text: value, class: '' };
        return `<span class="status-badge ${status.class}">${status.text}</span>`;
    }

    static actions(id, editCallback, deleteCallback) {
        return `
            <div class="table-actions">
                <button class="btn btn-primary btn-sm" onclick="event.stopPropagation(); ${editCallback}(${id})">编辑</button>
                <button class="btn btn-danger btn-sm" onclick="event.stopPropagation(); ${deleteCallback}(${id})">删除</button>
            </div>
        `;
    }

    static boolean(value) {
        return value ? '是' : '否';
    }

    static currency(value) {
        return `¥${parseFloat(value).toFixed(2)}`;
    }
}

// 导出到全局
window.DataTable = DataTable;
window.TableFormatters = TableFormatters;