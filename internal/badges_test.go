package internal

var expectedLabelDescriptionBadge = `
<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="108" height="20" role="img"
    aria-label="label: description">
    <title>label: description</title>
    <linearGradient id="s" x2="0" y2="100%">
        <stop offset="0" stop-color="#bbb" stop-opacity=".1" />
        <stop offset="1" stop-opacity=".1" />
    </linearGradient>
    <clipPath id="r">
        <rect width="108" height="20" rx="3" fill="#fff" />
    </clipPath>
    <g clip-path="url(#r)">
        <rect width="37" height="20" fill="#555" />
        <rect x="37" width="71" height="20" fill="#e05d44" />
        <rect width="108" height="20" fill="url(#s)" />
    </g>
    <g fill="#fff" text-anchor="middle" font-family="Verdana,Geneva,DejaVu Sans,sans-serif"
        text-rendering="geometricPrecision" font-size="110">
        <text aria-hidden="true" x="195" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)"
            textLength="270">label</text>
        <text x="195" y="140" transform="scale(.1)" fill="#fff" textLength="270">label</text>
        <text aria-hidden="true" x="715" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)"
            textLength="610">description</text>
        <text x="715" y="140" transform="scale(.1)" fill="#fff" textLength="610">description</text>
    </g>
</svg>
`
