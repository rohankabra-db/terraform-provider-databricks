package workspace

import (
	"context"

	"github.com/databricks/terraform-provider-databricks/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// DataSourceNotebook ...
func DataSourceNotebook() *schema.Resource {
	s := map[string]*schema.Schema{
		"path": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"format": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"DBC",
				"SOURCE",
				"HTML",
			}, false),
		},
		"content": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"language": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"object_type": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"object_id": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
	}
	return &schema.Resource{
		Schema: s,
		ReadContext: func(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
			notebooksAPI := NewNotebooksAPI(ctx, m)
			path := d.Get("path").(string)
			format := d.Get("format").(string)
			notebookContent, err := notebooksAPI.Export(path, format)
			if err != nil {
				return diag.FromErr(err)
			}
			d.SetId(path)
			// nolint
			d.Set("content", notebookContent)
			objectStatus, err := notebooksAPI.Read(d.Id())
			if err != nil {
				return diag.FromErr(err)
			}
			err = common.StructToData(objectStatus, s, d)
			if err != nil {
				return diag.FromErr(err)
			}
			return nil
		},
	}
}
