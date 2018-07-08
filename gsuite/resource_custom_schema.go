package gsuite

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"

	directory "google.golang.org/api/admin/directory/v1"
	"log"
)

func resourceCustomSchema() *schema.Resource {
	return &schema.Resource{
		Create: customSchemaCreate,
		Read:   customSchemaRead,
		Update: customSchemaUpdate,
		Delete: customSchemaDelete,
		Importer: &schema.ResourceImporter{
			State: customSchemaImporter,
		},

		Schema: map[string]*schema.Schema{

			"schema_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"etag": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"kind": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},


			"field": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"field_name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"field_type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"multi_valued": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},

						"read_access_type": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},

						"indexed": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},

						"numeric_indexing_spec": &schema.Schema{
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{

								Schema: map[string]*schema.Schema{
									"min_value": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
									},

									"max_value": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},

						"etag": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"kind": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"field_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// TODO: schemaCreate
func customSchemaCreate(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)

	customSchema := &directory.Schema{
		SchemaName: d.Get("schema_name").(string),
	}

	if _, ok := d.GetOk("fields"); ok {

		schemaFieldSpecs := []*directory.SchemaFieldSpec{}

		for i := 0; i < d.Get("fields.#").(int); i++ {
			fieldConfig := d.Get(fmt.Sprintf("posix_accounts.%d", i)).(map[string]interface{})
			fieldSpec := &directory.SchemaFieldSpec{}

			if fieldConfig["field_name"] != "" {
				log.Printf("[DEBUG] Setting posix %d gecos: %s", i, fieldConfig["field_name"].(string))
				fieldSpec.FieldId = fieldConfig["field_name"].(string)
			}

			if fieldConfig["field_type"] != "" {
				log.Printf("[DEBUG] Setting posix %d gecos: %s", i, fieldConfig["field_type"].(string))
				fieldSpec.FieldId = fieldConfig["field_type"].(string)
			}

			if fieldConfig["multi_valued"] != "" {
				log.Printf("[DEBUG] Setting posix %d gecos: %s", i, fieldConfig["multi_valued"].(string))
				fieldSpec.FieldId = fieldConfig["multi_valued"].(string)
			}

			if fieldConfig["read_access_type"] != "" {
				log.Printf("[DEBUG] Setting posix %d gecos: %s", i, fieldConfig["read_access_type"].(string))
				fieldSpec.FieldId = fieldConfig["read_access_type"].(string)
			}

			if fieldConfig["indexed"] != "" {
				log.Printf("[DEBUG] Setting posix %d gecos: %s", i, fieldConfig["indexed"].(string))
				fieldSpec.FieldId = fieldConfig["indexed"].(string)
			}

			if fieldConfig["numeric_indexing_spec"] != "" {

			}
			schemaFieldSpecs = append(schemaFieldSpecs, fieldSpec)
		}
		customSchema.Fields = schemaFieldSpecs
	}

	var createdSchema *directory.Schema
	var err error
	err = retry(func() error {
		createdSchema, err = config.directory.Schemas.Insert(customer_id, customSchema).Do()
		return err
	})

	if err != nil {
		return fmt.Errorf("Error creating user: %s", err)
	}

	d.SetId(createdSchema.SchemaId)
	log.Printf("[INFO] Created Schema: %s", createdSchema.SchemaName)
	return resourceUserRead(d, meta)
}


func customSchemaRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	var customSchema *directory.Schema
	var err error
	err = retry(func() error {
		customSchema, err = config.directory.Schemas.Get(d.Id(), d.Get("schema_id").(string)).Do()
		return err
	})

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Group %q", d.Get("name").(string)))
	}

	d.SetId(customSchema.SchemaId)

	d.Set("schema_name", customSchema.SchemaName)
	d.Set("etag", customSchema.Etag)
	d.Set("kind", customSchema.Kind)

	var fields []map[string]interface{}

	for _, field := range customSchema.Fields {
		var f map[string]interface{}

		f["etag"] = field.Etag
		f["field_id"] = field.FieldId
		f["field_name"] = field.FieldName
		f["field_type"] = field.FieldId
		f["multi_valued"] = field.MultiValued
		f["read_access_type"] = field.ReadAccessType
		f["indexed"] = field.Indexed
		f["kind"] = field.Kind

		if field.NumericIndexingSpec != nil {
			var numericIndexingSpec map[string]interface{}
			numericIndexingSpec["MinValue"] = field.NumericIndexingSpec.MinValue
			numericIndexingSpec["MaxValue"] = field.NumericIndexingSpec.MaxValue
			f["numericIndexingSpec"] = numericIndexingSpec
		}

		fields = append(fields, f)
	}

	d.Set("fields", fields)

	return nil
}

// TODO: schemaUpdate
func customSchemaUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func customSchemaDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	var err error
	err = retry(func() error {
		err = config.directory.Schemas.Delete(customer_id, d.Get("schema_id").(string)).Do()
		return err
	})
	if err != nil {
		return fmt.Errorf("Error deleting Schema: %s", err)
	}

	d.SetId("")
	return nil
}

// TODO: schemaImporter
func customSchemaImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	customSchema, err := config.directory.Schemas.Get(customer_id, d.Get("schema_id").(string)).Do()

	if err != nil {
		return nil, fmt.Errorf("Error fetching Schema. Make sure the schema exists: %s ", err)
	}

	d.SetId(customSchema.SchemaId)
	// TODO: Import Schema Properties

	return []*schema.ResourceData{d}, nil
}


