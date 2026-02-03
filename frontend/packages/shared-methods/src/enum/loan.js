export const tinOptions = [
	{ label: 'SSN', value: 'SSN' },
	{ label: 'ITIN', value: 'ITIN' },
	{ label: 'EIN', value: 'EIN' },
];
export const ethnicityOptions = [
	{ label: 'Hispanic or Latino', value: 'Hispanic or Latino' },
	{ label: 'Mexican', value: 'Mexican' },
	{ label: 'Cuban', value: 'Cuban' },
	{
		label: 'Other Hispanic or Latino',
		value: 'Other Hispanic or Latino',
		extras: {
			paramsKey: 'other_hispanic_or_latino_name',
			eventType: 'input',
			placeholder: 'Enter ethnicity, for example: Colombian, Dominican, Nicaraguan and so on',
		},
	},
	{ label: 'Not Hispanic or Latino', value: 'Not Hispanic or Latino' },
	{
		label: 'I do not wish to provide this information',
		value: 'Hidden',
		isRadio: true,
	},
];

export const raceOptions = [
	{
		label: 'American Indian or Alaska Native',
		value: 'American Indian or Alaska Native',
		extras: {
			paramsKey: 'american_indian_or_alaska_native_name',
			eventType: 'input',
			placeholder: 'Name of enrolled or principal tribe',
		},
	},
	{
		label: 'Asian',
		value: 'Asian',
		extras: {
			eventType: 'checkbox',
			paramsKey: 'asian_values',
			options: [
				{
					label: 'Asian Indian',
					value: 'Asian Indian',
				},
				{
					label: 'Chinese',
					value: 'Chinese',
				},
				{
					label: 'Filipino',
					value: 'Filipino',
				},
				{
					label: 'Japanese',
					value: 'Japanese',
				},
				{
					label: 'Korean',
					value: 'Korean',
				},
				{
					label: 'Vietnamese',
					value: 'Vietnamese',
				},
				{
					label: 'Other Asian',
					value: 'Other Asian',
					extras: {
						paramsKey: 'asian_other_name',
						eventType: 'input',
						placeholder: 'Enter race, For example: Hmong, Lootian, Thai and so on',
					},
				},
			],
		},
	},
	{ label: 'Black or African American', value: 'Black or African American' },
	{
		label: 'Native Hawaiian or Other Pacific Islander',
		value: 'Native Hawaiian or Other Pacific Islander',
	},
	{
		label: 'Native Hawaiian Guamanian or Chamorro Samoan',
		value: 'Native Hawaiian Guamanian or Chamorro Samoan',
	},
	{
		label: 'Other Pacific Islander',
		value: 'Other Pacific Islander',
		extras: {
			paramsKey: 'other_pacific_islander_name',
			eventType: 'input',
			placeholder: 'Enter race, For example: Fijian, Tongan and so on',
		},
	},
	{ label: 'White', value: 'White' },
	{
		label: 'I do not wish to provide this information',
		value: 'Hidden',
		isRadio: true,
	},
];
