package drone

import (
	"context"
	"fmt"

	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mavimo/terraform-provider-drone/drone/utils"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Description: `Manage a user.

		~> In order to use the _drone_user_ resource you must have admin privileges within your Drone environment.`,
		Schema: map[string]*schema.Schema{
			"login": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Login name",
			},
			"active": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Is the user active?",
			},
			"admin": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Is the user an admin?",
			},
			"machine": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Is the user a machine?",
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user's access token",
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Exists:        resourceUserExists,
	}
}

func resourceUserCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(drone.Client)

	user, err := client.UserCreate(createUser(data))

	if err != nil {
		return diag.Errorf("Unable to create user %s", user.Login)
	}

	return readUser(data, user, err)
}

func resourceUserUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(drone.Client)

	user, err := client.User(data.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	user, err = client.UserUpdate(user.Login, updateUser(data))

	return readUser(data, user, err)
}

func resourceUserRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(drone.Client)

	user, err := client.User(data.Id())
	if err != nil {
		return diag.Errorf("failed to read Drone user with id: %s", data.Id())
	}

	return readUser(data, user, err)
}

func resourceUserDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(drone.Client)

	err := client.UserDelete(data.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceUserExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(drone.Client)

	login := data.Id()

	user, err := client.User(login)
	if err != nil {
		return false, fmt.Errorf("failed to read Drone user with id: %s", data.Id())
	}

	exists := user.Login == login

	return exists, err
}

func createUser(data *schema.ResourceData) (user *drone.User) {
	user = &drone.User{
		Login:   data.Get("login").(string),
		Active:  data.Get("active").(bool),
		Admin:   data.Get("admin").(bool),
		Machine: data.Get("machine").(bool),
	}

	return user
}

func updateUser(data *schema.ResourceData) (user *drone.UserPatch) {
	userPatch := &drone.UserPatch{
		Active:  utils.Bool(data.Get("active").(bool)),
		Admin:   utils.Bool(data.Get("admin").(bool)),
		Machine: utils.Bool(data.Get("machine").(bool)),
	}
	return userPatch
}

func readUser(data *schema.ResourceData, user *drone.User, err error) diag.Diagnostics {
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(user.Login)

	data.Set("login", user.Login)
	data.Set("active", user.Active)
	data.Set("machine", user.Machine)
	data.Set("admin", user.Admin)
	if user.Token != "" {
		data.Set("token", user.Token)
	}

	return diag.Diagnostics{}
}
