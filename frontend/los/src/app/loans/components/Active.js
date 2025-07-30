'use client';
import React, { useEffect, useMemo, useState } from 'react';
import FilterTable from '@/components/FilterTable';
import { activeColumns } from './column';
const ActiveLoans = () => {
    const [type, setType] = useState('');
    const [searchValue, setSearchValue] = useState('');
    const [loanPurpose, setLoanPurpose] = useState([]);
    const [loanType, setLoanType] = useState([]);
    const [lender, setLender] = useState([]);
    const [status, setStatus] = useState([]);
    const filterList = useMemo(
        () => [
            {
                type: 'select-input',
                selectValue: type,
                value: searchValue,
                selectProps: {
                    onChange: (e) => setType(e),
                },
                onChange: ({ target: { value } }) => setSearchValue(value),
                options: [
                    {
                        label: 'Address',
                        value: 'address',
                    },
                    {
                        label: 'State',
                        value: 'state',
                    },
                    {
                        label: 'City',
                        value: 'city',
                    },
                    {
                        label: 'County',
                        value: 'county',
                    },
                    {
                        label: 'Zip Code',
                        value: 'zipCode',
                    },
                ],
            },
            {
                type: 'treeSelect',
                label: 'Loan purpose',
                popupMatchSelectWidth: 172,
                placeholder: 'Loan purpose',
                value: loanPurpose,
                onChange: (value) => setLoanPurpose(value),
                treeData: [
                    {
                        title: 'Purchase',
                        value: 'purchase',
                        key: 'purchase',
                        isLeaf: true,
                    },
                    {
                        title: 'Rate-term refi',
                        value: 'rateTermRefi',
                        key: 'rateTermRefi',
                        isLeaf: true,
                    },
                    {
                        title: 'Cash-out refi',
                        value: 'cashOutRefi',
                        key: 'cashOutRefi',
                        isLeaf: true,
                    },
                ],
            },
            {
                type: 'treeSelect',
                label: 'Loan type',
                popupMatchSelectWidth: 195,
                value: loanType,
                placeholder: 'Loan type',
                onChange: (value) => {
                    console.log('change',value);
                    
                    setLoanType(value)
                },
                treeData: [
                    {
                        title: 'Bridge',
                        value: 'bridge',
                        key: 'bridge',
                        isLeaf: true,
                    },
                    {
                        title: 'Rehab',
                        value: 'rehab',
                        key: 'rehab',
                        isLeaf: true,
                    },
                    {
                        title: 'Ground up construction',
                        value: 'groundUpConstruction',
                        key: 'groundUpConstruction',
                        isLeaf: true,
                    },
                ],
            },
            {
                type: 'treeSelect',
                label: 'Lender',
                placeholder: 'Lender',
                popupMatchSelectWidth: 172,
                value: lender,
                onChange: (value) => setLender(value),
                treeData: [
                    {
                        title: 'Bridge',
                        value: 'bridge',
                        key: 'bridge',
                        isLeaf: true,
                    },
                    {
                        title: 'Rehab',
                        value: 'rehab',
                        key: 'rehab',
                        isLeaf: true,
                    },
                    {
                        title: 'Ground up construction',
                        value: 'groundUpConstruction',
                        key: 'groundUpConstruction',
                        isLeaf: true,
                    },
                ],
            },
            
        ],
        [
            searchValue,
            type,
            loanPurpose,
            loanType,
            lender,
            status,
        ],
    );

    return (
        <FilterTable
            columns={activeColumns}
            scroll={{ x: 3500, y: 'calc(100% - 54px)' }}
            dataSource={[{ status: 'quoting', id: 1 }]}
            filterList={filterList}
        />
    );
};
export default ActiveLoans;
